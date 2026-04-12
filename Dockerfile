FROM golang:1.25.6 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./cmd/bin/app ./cmd

FROM alpine:3.19

RUN apk add --no-cache tzdata curl

WORKDIR /app

COPY --from=builder /app/cmd/bin/app .

HEALTHCHECK --interval=10s --timeout=10s --retries=3 \
    CMD curl -fs http://localhost:8080/health || exit 1

CMD ["./app"]
