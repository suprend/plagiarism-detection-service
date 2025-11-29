# syntax = docker/dockerfile:1

ARG GO_VERSION=1.25
ARG ALPINE_VERSION=3.20
ARG TARGETARCH

FROM golang:${GO_VERSION}-alpine AS builder

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=${TARGETARCH:-amd64} \
    go build -trimpath -ldflags="-s -w" -o /app/bin/server ./cmd/server

FROM alpine:${ALPINE_VERSION}

RUN apk --no-cache add ca-certificates && addgroup -S app && adduser -S app -G app

WORKDIR /app

COPY --from=builder /app/bin/server .
COPY openapi.yaml .

ENV PORT=8082 \
    FILESTORAGE_URL=http://filestorage:8080 \
    PLAGIARISM_URL=http://plagiarism:8081

EXPOSE 8082

USER app

ENTRYPOINT ["./server"]
