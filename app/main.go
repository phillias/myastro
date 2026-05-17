package main

import (
	"astrology-app/jyotish"
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/getsentry/sentry-go"
)

//go:embed static/*
var staticFiles embed.FS

var sentryDSN = "" // set via -ldflags or SENTRY_DSN (see AGENTS.md)

func main() {
	if dsn := sentryDSN; dsn != "" {
		if err := sentry.Init(sentry.ClientOptions{
			Dsn:              dsn,
			TracesSampleRate: 0.1,
		}); err != nil {
			log.Printf("Warning: Sentry init failed: %v", err)
		} else {
			log.Println("Sentry initialized")
		}
	}

	dist, _ := fs.Sub(staticFiles, "static")
	fileServer := http.FileServer(http.FS(dist))

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fileServer.ServeHTTP(w, r)
	})

	mux.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	mux.HandleFunc("/api/calculate", func(w http.ResponseWriter, r *http.Request) {
		// Recover from panics
		defer func() {
			if err := recover(); err != nil {
				sentry.CurrentHub().Recover(err)
				http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
			}
		}()

		latParam := r.URL.Query().Get("lat")
		lonParam := r.URL.Query().Get("lon")
		date := r.URL.Query().Get("date")
		timeStr := r.URL.Query().Get("time")

		log.Printf("API Request: lat=%s, lon=%s, date=%s, time=%s", latParam, lonParam, date, timeStr)

		// Input validation
		if latParam == "" || lonParam == "" || date == "" || timeStr == "" {
			jsonError(w, "lat, lon, date, and time are required", http.StatusBadRequest)
			return
		}

		lat, err := strconv.ParseFloat(latParam, 64)
		if err != nil || lat < -90 || lat > 90 {
			jsonError(w, "lat must be a float between -90 and 90", http.StatusBadRequest)
			return
		}

		lon, err := strconv.ParseFloat(lonParam, 64)
		if err != nil || lon < -180 || lon > 180 {
			jsonError(w, "lon must be a float between -180 and 180", http.StatusBadRequest)
			return
		}

		birthTime, err := time.Parse("2006-01-02 15:04", fmt.Sprintf("%s %s", date, timeStr))
		if err != nil {
			jsonError(w, "date must be YYYY-MM-DD, time must be HH:MM", http.StatusBadRequest)
			return
		}

		// Check cache first
		cacheKey := jyotish.CacheKey(lat, lon, date, timeStr)
		if cached, hit := jyotish.CachedResponse(cacheKey); hit {
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("X-Cache", "HIT")
			w.Write(cached)
			return
		}

		// Resolve timezone from lat/lon
		loc := jyotish.GetTimeZone(lat, lon)
		birthLocal := birthTime.In(loc)

		// Full chart calculation (all planets, lagna, houses) via Swiss Ephemeris
		chart := jyotish.CalculateChart(birthLocal, lat, lon)

		var sunLon, moonLon float64
		for _, p := range chart.Planets {
			switch p.Name {
			case "Sun":
				sunLon = p.Longitude
			case "Moon":
				moonLon = p.Longitude
			}
		}
		log.Printf("Chart for %v: Lagna=%.2f° %s, Sun=%.2f°, Moon=%.2f°",
			birthLocal, chart.Lagna.Degree, chart.Lagna.Sign, sunLon, moonLon)

		dashaResult := jyotish.CalculateDashas(moonLon, birthLocal)

		data := jyotish.JyotishData{
			Sign:          chart.Lagna.Sign,
			Nakshatra:     jyotish.NakshatraFromLongitude(moonLon),
			NakshatraPada: jyotish.NakshatraPada(moonLon),
			LaganaDeg:     chart.Lagna.Longitude,
			Panchanga:     jyotish.CalculatePanchangam(sunLon, moonLon),
			Remedy:        jyotish.GetRemedy(chart.Lagna.Sign),
			Chart:         &chart,
			Dasha:         &dashaResult,
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Cache", "MISS")
		json.NewEncoder(w).Encode(data)

		body, _ := json.Marshal(data)
		jyotish.CacheResponse(cacheKey, body)
	})

	mux.HandleFunc("/api/compatibility", func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				sentry.CurrentHub().Recover(err)
				http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
			}
		}()

		lat1 := r.URL.Query().Get("lat1")
		lon1 := r.URL.Query().Get("lon1")
		date1 := r.URL.Query().Get("date1")
		time1 := r.URL.Query().Get("time1")
		lat2 := r.URL.Query().Get("lat2")
		lon2 := r.URL.Query().Get("lon2")
		date2 := r.URL.Query().Get("date2")
		time2 := r.URL.Query().Get("time2")

		if lat1 == "" || lon1 == "" || date1 == "" || time1 == "" ||
			lat2 == "" || lon2 == "" || date2 == "" || time2 == "" {
			jsonError(w, "lat1, lon1, date1, time1, lat2, lon2, date2, time2 are required", http.StatusBadRequest)
			return
		}

		la1, err := strconv.ParseFloat(lat1, 64)
		if err != nil || la1 < -90 || la1 > 90 {
			jsonError(w, "lat1 must be a float between -90 and 90", http.StatusBadRequest)
			return
		}
		lo1, err := strconv.ParseFloat(lon1, 64)
		if err != nil || lo1 < -180 || lo1 > 180 {
			jsonError(w, "lon1 must be a float between -180 and 180", http.StatusBadRequest)
			return
		}
		la2, err := strconv.ParseFloat(lat2, 64)
		if err != nil || la2 < -90 || la2 > 90 {
			jsonError(w, "lat2 must be a float between -90 and 90", http.StatusBadRequest)
			return
		}
		lo2, err := strconv.ParseFloat(lon2, 64)
		if err != nil || lo2 < -180 || lo2 > 180 {
			jsonError(w, "lon2 must be a float between -180 and 180", http.StatusBadRequest)
			return
		}

		t1, err := time.Parse("2006-01-02 15:04", fmt.Sprintf("%s %s", date1, time1))
		if err != nil {
			jsonError(w, "date1 must be YYYY-MM-DD, time1 must be HH:MM", http.StatusBadRequest)
			return
		}
		t2, err := time.Parse("2006-01-02 15:04", fmt.Sprintf("%s %s", date2, time2))
		if err != nil {
			jsonError(w, "date2 must be YYYY-MM-DD, time2 must be HH:MM", http.StatusBadRequest)
			return
		}

		loc1 := jyotish.GetTimeZone(la1, lo1)
		loc2 := jyotish.GetTimeZone(la2, lo2)
		birthLocal1 := t1.In(loc1)
		birthLocal2 := t2.In(loc2)

		chartA := jyotish.CalculateChart(birthLocal1, la1, lo1)
		chartB := jyotish.CalculateChart(birthLocal2, la2, lo2)

		result := jyotish.CalculateSynastry(chartA, chartB)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	})

	mux.HandleFunc("/api/dasha", func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				sentry.CurrentHub().Recover(err)
				http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
			}
		}()

		latParam := r.URL.Query().Get("lat")
		lonParam := r.URL.Query().Get("lon")
		date := r.URL.Query().Get("date")
		timeStr := r.URL.Query().Get("time")

		if latParam == "" || lonParam == "" || date == "" || timeStr == "" {
			jsonError(w, "lat, lon, date, and time are required", http.StatusBadRequest)
			return
		}

		lat, err := strconv.ParseFloat(latParam, 64)
		if err != nil || lat < -90 || lat > 90 {
			jsonError(w, "lat must be a float between -90 and 90", http.StatusBadRequest)
			return
		}
		lon, err := strconv.ParseFloat(lonParam, 64)
		if err != nil || lon < -180 || lon > 180 {
			jsonError(w, "lon must be a float between -180 and 180", http.StatusBadRequest)
			return
		}

		birthTime, err := time.Parse("2006-01-02 15:04", fmt.Sprintf("%s %s", date, timeStr))
		if err != nil {
			jsonError(w, "date must be YYYY-MM-DD, time must be HH:MM", http.StatusBadRequest)
			return
		}

		loc := jyotish.GetTimeZone(lat, lon)
		birthLocal := birthTime.In(loc)

		chart := jyotish.CalculateChart(birthLocal, lat, lon)

		var moonLon float64
		for _, p := range chart.Planets {
			if p.Name == "Moon" {
				moonLon = p.Longitude
				break
			}
		}

		result := jyotish.CalculateDashas(moonLon, birthLocal)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	})

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func jsonError(w http.ResponseWriter, msg string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}
