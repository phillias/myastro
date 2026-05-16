# Stage 1: Build Frontend
FROM node:20-bookworm AS frontend-build
WORKDIR /app
COPY myastrology/frontend/package*.json ./
RUN npm install
COPY myastrology/frontend/ .
RUN npm run build

# Stage 2: Build Backend
FROM golang:1.23-bookworm AS backend-build
RUN apt-get update && apt-get install -y build-essential libswe-dev
WORKDIR /app
COPY myastrology/app/go.* ./
RUN go mod download
COPY myastrology/app/ .
# Embed frontend build
COPY --from=frontend-build /app/dist ./static

# Set CGO flags to point to headers and libs
# -I /usr/include finds the .h files
# -L /usr/lib/x86_64-linux-gnu links the binary .so files and math library
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
CMD ["./astrology-app"]
