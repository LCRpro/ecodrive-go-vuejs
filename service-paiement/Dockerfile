FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . .
RUN go build -o app .
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/app .
COPY .env .
EXPOSE 8004
CMD ["./app"]