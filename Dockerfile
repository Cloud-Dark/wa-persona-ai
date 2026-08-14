FROM golang:1.22-alpine AS builder

RUN apk add --no-cache gcc musl-dev sqlite-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-s -w" -o /app/bin/wa-persona-ai ./cmd/wa-persona-ai

# ─── Runtime ─────────────────────────────────────────────────────────
FROM alpine:3.19

RUN apk add --no-cache ca-certificates sqlite-libs tzdata

WORKDIR /app

COPY --from=builder /app/bin/wa-persona-ai .
COPY persona/ ./persona/
COPY config/config.example.yaml ./config/config.yaml

RUN mkdir -p data/sessions data/memory

ENV WPA_CONFIG_PATH=/app/config/config.yaml

VOLUME ["/app/data", "/app/config", "/app/persona"]

ENTRYPOINT ["./wa-persona-ai"]
