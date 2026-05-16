# MyAstro — Vedic Astrology (Jyotish) PWA

A full-stack Progressive Web Application for Vedic astrology calculations — computes your zodiac sign, nakshatra, panchanga (Hindu calendar), and spiritual remedies based on birth time and location.

## Overview

MyAstro uses the **Swiss Ephemeris** library for accurate astronomical calculations and the **Lahiri Ayanamsa** (precession correction) to convert tropical coordinates to the sidereal (Vedic) system.

## Architecture

```
myastro/
├── app/                          # Go backend API
│   ├── main.go                   # HTTP server (embed static, /api/calculate endpoint)
│   ├── go.mod                    # Go module
│   └── jyotish/                  # Astrological calculation engine
│       ├── calculator.go         # Zodiac sign & nakshatra from sidereal longitude
│       ├── ephemeris.go          # Swiss Ephemeris Cgo bindings (Sun/Moon positions)
│       ├── panchangam.go         # Tithi, Nakshatra, Yoga, Karana computation
│       └── remedies.go           # Sign-to-mantra/yantra/color mappings
├── frontend/                     # React + Vite + TypeScript SPA
│   ├── src/
│   │   ├── App.tsx               # Main app (map picker, date/time inputs, results dialog)
│   │   ├── MapPicker.tsx         # Leaflet map with geosearch
│   │   ├── leaflet-fix.ts        # Leaflet icon fix for bundlers
│   │   └── index.css             # Tailwind + Leaflet styles
│   ├── public/manifest.json      # PWA manifest
│   └── package.json
├── Dockerfile                    # Multi-stage Docker build
├── compose.yml                   # Docker Compose service definition
└── AGENTS.md                     # Development guide for AI agents
```

### Backend — Go API (`:8080`)

| Endpoint | Method | Parameters | Returns |
|----------|--------|------------|---------|
| `/`       | GET    | — | Serves embedded frontend static files |
| `/api/calculate` | GET | `lat`, `lon`, `date` (YYYY-MM-DD), `time` (HH:MM) | JSON with sign, nakshatra, lagana_degree, remedy (mantra/yantra/color), panchanga (tithi, nakshatra, yoga, karana) |

**Astrological calculations:**
- **Sidereal conversion** — Applies Lahiri Ayanamsa (~24.16° precession for 2026)
- **Rashi (zodiac sign)** — Determined by sidereal longitude (30° divisions)
- **Nakshatra (lunar mansion)** — 27 divisions of the ecliptic
- **Panchanga** — Five-limbed Hindu calendar: Tithi (lunar day), Nakshatra, Yoga, Karana
- **Remedies** — Vedic mantras, yantras, and colors per sign

The backend embeds the frontend build output using Go's `//go:embed` directive, resulting in a single deployable binary.

### Frontend — React SPA

- **Interactive map** — Leaflet with OpenStreetMap and geosearch for birth location
- **Date/time pickers** — Native HTML5 inputs
- **Results display** — Radix UI tabs (Natal view / Panchanga view) in a modal dialog
- **PWA** — Manifest configured for standalone install
- **Styling** — Tailwind CSS, dark theme (`slate-900`), Radix UI primitives, Lucide icons

## Quick Start

### Prerequisites

- Go 1.23+
- Node.js 20+
- Swiss Ephemeris library (`libswe-dev`)
- Docker (optional, for containerized build)

### Development

**1. Frontend dev server (hot reload):**
```bash
cd frontend && npm install && npm run dev
```

**2. Backend (serve API + embedded frontend):**
```bash
cd app && go run .
```

**3. Open http://localhost:8080**

### Production (Docker)

```bash
docker compose up --build
```

This runs a multi-stage build:
1. Builds the React frontend with Vite
2. Builds the Go binary with embedded frontend (requires `libswe-dev`)
3. Creates a slim Debian runtime image with just `libswe` and the binary

### API Usage

```bash
curl 'http://localhost:8080/api/calculate?lat=40.7128&lon=-74.0060&date=2026-05-15&time=14:30'
```

Example response:
```json
{
  "sign": "Taurus",
  "nakshatra": "Mrigashira",
  "lagana_degree": 45.84,
  "remedy": {
    "mantra": "Om Dram Drim Draum Sah Shukraya Namah",
    "yantra": "Shukra Yantra",
    "color": "White"
  },
  "panchanga": {
    "tithi": "Chaturthi",
    "nakshatra": "Mrigashira",
    "yoga": "Yoga 3",
    "karana": "Karana 14"
  }
}
```

## Dependencies

### Backend
- **Swiss Ephemeris** (`libswe`) — High-precision planetary ephemeris via Cgo bindings
- Go standard library only (no external Go modules)

### Frontend
- **React 19** + **TypeScript**
- **Vite 6** — Build tool
- **Tailwind CSS 3** — Utility-first styling
- **Leaflet** + **react-leaflet** — Interactive maps
- **leaflet-geosearch** — Location search
- **Radix UI** — Accessible dialog and tabs primitives
- **Lucide React** — Icons

## License

MIT
