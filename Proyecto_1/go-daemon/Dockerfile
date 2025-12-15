# Etapa de build
FROM golang:1.24 AS builder
WORKDIR /app
COPY . .
RUN go build -o go-daemon ./cmd/main.go

# Imagen final
FROM debian:stable-slim
WORKDIR /app
COPY --from=builder /app/go-daemon /app/go-daemon
COPY config.yaml /app/
CMD ["./go-daemon"]
