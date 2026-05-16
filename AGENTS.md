# myastrology Repository Agent Guide

## Project Structure
- **Frontend**: React/Vite app in `/frontend`
- **Backend**: Go API server in `/app` serving embedded frontend
- **API Endpoint**: `/api/calculate` (GET) returns astrological calculations
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
- Go version: 1.23 (see go.mod)
- Node version: 20 (see Dockerfile)

### API Usage
Endpoint: `GET /api/calculate`
Parameters:
- `lat`: latitude (float)
- `lon`: longitude (float)  
- `date`: YYYY-MM-DD format
- `time`: HH:MM format (24-hour)

Returns JSON with astrological calculations including zodiac sign, remedies, panchanga, etc.

### Testing
No test framework configured in visible config files. Manual verification via:
1. Frontend: Visual inspection in browser
2. Backend: curl/API testing of `/api/calculate` endpoint