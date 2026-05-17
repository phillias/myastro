# myastrology Repository Agent Guide

## Project Structure
- **Frontend**: React/Vite app in `/frontend`
- **Backend**: Go API server in `/app` serving embedded frontend
- **API Endpoints**:
  - `/api/calculate` (GET) — individual birth chart
  - `/api/compatibility` (GET) — synastry between two charts
- **Planet Computation**: `jyotish/planets.go` — all 9 Vedic planets (Sun–Saturn + Rahu/Ketu) via Swiss Ephemeris
- **Synastry Engine**: `jyotish/synastry.go` — cross-chart aspects, Vedic graha drishti, house overlays
- **Static Files**: Frontend build output embedded in Go binary via `go:embed`

## Development Commands

### Frontend
```bash
# Start dev server
cd frontend && npm run dev

# Build for production
cd frontend && npm run build

# Lint code
cd frontend && npm run lint

# Preview production build
cd frontend && npm run preview
```

### Backend
```bash
# Build and run (development)
go run .  # From /app directory

# Build binary
go build -o astrology-app  # From /app directory
```

### Full Application (Docker)
```bash
# Build and run with Docker
docker compose up --build

# Just start (if already built)
docker compose up
```

## Important Notes

### Build Order Matters
1. Frontend must be built before backend (Dockerfile handles this)
2. Backend embeds frontend static files at compile time
3. For local dev: run frontend dev server separately, backend serves API only

### Environment & Dependencies
- Backend requires `libswe-dev` for Cgo bindings to Swiss Ephemeris library
- Frontend uses Vite with React + TypeScript + Tailwind CSS
- Go version: 1.24 (see go.mod)
- Node version: 20 (see Dockerfile)
- Sentry Go SDK v0.29.0 (disabled unless SENTRY_DSN passed via -ldflags)
- TZF (ringsaturn/tzf) for timezone resolution from lat/lon

### API: Individual Chart `GET /api/calculate`
Parameters:
- `lat`, `lon`: latitude/longitude (float)
- `date`: YYYY-MM-DD
- `time`: HH:MM (24-hour)

Returns: `{ sign, nakshatra, nakshatra_pada, lagana_degree, remedy, panchanga }`

### API: Synastry `GET /api/compatibility`
Parameters (two persons):
- `lat1`, `lon1`, `date1`, `time1` — person A
- `lat2`, `lon2`, `date2`, `time2` — person B

Returns: `{ chart_a, chart_b, aspects[], vedic_drishti[], house_overlays[], summary }`

Each chart contains Lagna + 9 planets with sign, degree, nakshatra, pada, house, retrograde status.
Aspects: conjunction/sextile/square/trine/opposition with orb values.
Vedic drishti: Jupiter(5/7/9), Saturn(3/7/10), Mars(4/7/8) applied cross-chart.

### Testing
```bash
# Run all backend tests
cd app && go test ./jyotish/ -v -count=1
```
Tests package: calculator logic, panchanga edge cases, house position math,
angular distances, aspect detection with synthetic data, full chart structure validation.
23 tests across 6 test files.