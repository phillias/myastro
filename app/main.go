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

		data := jyotish.JyotishData{}

		sunLon, moonLon := jyotish.GetSunMoonLongitudes(birthLocal)
		log.Printf("Ephemeris for %v: SunLon=%.2f, MoonLon=%.2f", birthLocal, sunLon, moonLon)

		data.LaganaDeg = jyotish.GetLagna(birthLocal, lat, lon)
		data.Sign = jyotish.RashiFromLongitude(data.LaganaDeg)
		data.Nakshatra = jyotish.NakshatraFromLongitude(moonLon)
		data.NakshatraPada = jyotish.NakshatraPada(moonLon)
		data.Panchanga = jyotish.CalculatePanchangam(sunLon, moonLon)
		data.Remedy = jyotish.GetRemedy(data.Sign)

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Cache", "MISS")
		json.NewEncoder(w).Encode(data)

		body, _ := json.Marshal(data)
		jyotish.CacheResponse(cacheKey, body)
	})

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func jsonError(w http.ResponseWriter, msg string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}
