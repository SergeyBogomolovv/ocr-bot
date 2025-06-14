FROM golang:1.24.4-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o bot *.go

FROM debian:12-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
    tesseract-ocr \
    tesseract-ocr-rus \
    tesseract-ocr-eng \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

RUN useradd -m nonroot

WORKDIR /app

COPY --from=builder /app/bot ./bot

USER nonroot

ENTRYPOINT ["./bot"]
