# ---- Build stage ----
FROM golang:1.25-alpine AS build
WORKDIR /src

# Cache dependencies first
COPY go.mod go.sum ./
RUN go mod download

# Copy source and build static binary
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/server .

# ---- Runtime stage (no package manager needed) ----
FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /app/backend

# The Go app resolves env via godotenv.Load("../containerization/.env")
# relative to its working directory. The compose file mounts the .env there
# and provides /data via a named volume.
ENV TORRENT_DATA_DIR=/data

COPY --from=build /bin/server /app/backend/server
EXPOSE 8080
ENTRYPOINT ["/app/backend/server"]
