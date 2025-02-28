# Используем минимальный образ Golang для компиляции
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o app cmd/main.go

FROM alpine:latest
COPY --from=builder /app/app /app
CMD ["/app"]