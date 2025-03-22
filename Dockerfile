FROM golang:1.23.6-alpine AS builder
WORKDIR /app
COPY . .
RUN go version
RUN go build -o song-library ./cmd/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/song-library .
COPY .env .
COPY internal/migration/migrations ./migrations
EXPOSE 8080
CMD ["./song-library"]