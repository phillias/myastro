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
)

//go:embed static/*
var staticFiles embed.FS

func main() {
	dist, _ := fs.Sub(staticFiles, "static")
	fileServer := http.FileServer(http.FS(dist))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fileServer.ServeHTTP(w, r)
	})

	http.HandleFunc("/api/calculate", func(w http.ResponseWriter, r *http.Request) {
		latParam := r.URL.Query().Get("lat")
		lonParam := r.URL.Query().Get("lon")
		date := r.URL.Query().Get("date") // Expects YYYY-MM-DD
		timeStr := r.URL.Query().Get("time") // Expects HH:MM

		log.Printf("API Request: lat=%s, lon=%s, date=%s, time=%s", latParam, lonParam, date, timeStr)
		
		// Parse birth time
		birthTime, err := time.Parse("2006-01-02 15:04", fmt.Sprintf("%s %s", date, timeStr))
		if err != nil {
			log.Printf("Error parsing date/time: %v, defaulting to Now", err)
			birthTime = time.Now()
		}

		lon, _ := strconv.ParseFloat(lonParam, 64)

		// 1. Calculate Natal Data
		data := jyotish.Calculate(lon)
		data.Remedy = jyotish.GetRemedy(data.Sign)
		
		// 2. Real-time Ephemeris calculation based on birthTime
		sunLon, moonLon := jyotish.GetSunMoonLongitudes(birthTime)
		log.Printf("Ephemeris for %v: SunLon=%.2f, MoonLon=%.2f", birthTime, sunLon, moonLon)
		
		// 3. Panchanga based on real ephemeris
		data.Panchanga = jyotish.CalculatePanchangam(sunLon, moonLon)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(data)
	})

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
