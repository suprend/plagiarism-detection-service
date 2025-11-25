# syntax = docker/dockerfile:1

ARG GO_VERSION=1.25
ARG ALPINE_VERSION=3.20
ARG TARGETARCH
ARG VCS_REF=unknown
ARG BUILD_DATE=unknown

# ---- build stage: compile Go binary with cached deps
FROM golang:${GO_VERSION}-alpine AS builder

WORKDIR /app

# Copy module files separately to maximize cache hits for dependencies
COPY go.mod go.sum ./
RUN go mod download

# Now copy the remaining source (invalidates cache only when needed)
COPY . .

# Produce a small, reproducible static binary
RUN CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=${TARGETARCH:-amd64} \
    go build -trimpath -ldflags="-s -w" -o /app/bin/server ./cmd/server

# ---- runtime stage: minimal image, non-root user, metadata labels
FROM alpine:${ALPINE_VERSION}

# re-declare build args for this stage (required by Docker)
ARG VCS_REF=unknown
ARG BUILD_DATE=unknown

LABEL org.opencontainers.image.title="filestorage" \
      org.opencontainers.image.description="file storage microservice for plagiarism detection platform" \
      org.opencontainers.image.revision="${VCS_REF}" \
      org.opencontainers.image.created="${BUILD_DATE}"

# Install only required runtime packages
RUN apk --no-cache add ca-certificates

# Drop privileges for better security
RUN addgroup -S app && adduser -S app -G app

WORKDIR /app

COPY --from=builder /app/bin/server .

USER app

# Documentation port; actual port still configurable via env
EXPOSE 8080

ENTRYPOINT ["./server"]

