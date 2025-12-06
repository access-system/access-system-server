FROM golang:1.24.6 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o access-system-server cmd/main.go


FROM alpine:3.20.1
WORKDIR /app

COPY --from=builder /app/access-system-server .
COPY --from=builder /app/.env .env

EXPOSE 8081
CMD ["./access-system-server"]
