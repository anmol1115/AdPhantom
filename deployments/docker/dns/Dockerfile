FROM golang:1.25.4-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

# build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 \
    go build -o dns-adblocker ./cmd/dns

# ---- runtime ----
FROM alpine:3.23.2

RUN adduser -D dnsuser
USER dnsuser
WORKDIR /app

COPY --from=builder /app/dns-adblocker .

EXPOSE 53/udp
EXPOSE 53/tcp

CMD ["./dns-adblocker"]
