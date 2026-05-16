# Stage 1: Build Frontend
FROM node:20-bookworm AS frontend-build
WORKDIR /app
COPY frontend/package*.json ./
RUN npm install
COPY frontend/ .
RUN npm run build

# Stage 2: Build Backend
FROM golang:1.24-bookworm AS backend-build
RUN apt-get update && apt-get install -y build-essential libswe-dev
WORKDIR /app
COPY app/go.* ./
RUN go mod download
COPY app/ .
COPY --from=frontend-build /app/dist ./static
ENV CGO_CFLAGS="-I/usr/include"
ENV CGO_LDFLAGS="-L/usr/lib/x86_64-linux-gnu -lswe -lm"
RUN CGO_ENABLED=1 go build -o astrology-app .

# Final Stage: Run
FROM debian:bookworm-slim
RUN apt-get update && apt-get install -y libswe2.0 && rm -rf /var/lib/apt/lists/*
WORKDIR /root/
COPY --from=backend-build /app/astrology-app .
COPY --from=backend-build /app/static ./static
EXPOSE 8080
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:8080/api/health || exit 1
CMD ["./astrology-app"]
