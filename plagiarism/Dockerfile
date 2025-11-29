# syntax=docker/dockerfile:1
ARG GO_VERSION=1.25

FROM golang:${GO_VERSION} AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/bin/plagiarism ./cmd/server

FROM debian:bookworm-slim

WORKDIR /app

RUN mkdir -p /app/plagiarism/reports && useradd -u 10001 appuser

COPY --from=builder /app/bin/plagiarism /app/server
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

ENV PORT=8080 \
    FILESTORAGE_URL=http://filestorage:8080 \
    MATCH_THRESHOLD=0.8 \
    WORKER_COUNT=1

EXPOSE 8080

ENTRYPOINT ["/entrypoint.sh"]
