FROM golang:1.25-alpine AS build
WORKDIR /app
COPY . .
RUN go build -o wordcloud ./cmd/server

FROM alpine:3.19
WORKDIR /app
COPY --from=build /app/wordcloud /usr/local/bin/wordcloud
EXPOSE 8083
ENTRYPOINT ["/usr/local/bin/wordcloud"]
