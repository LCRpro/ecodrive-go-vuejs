FROM golang:1.22-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o service-user .

FROM alpine:3.20

WORKDIR /app

COPY --from=build /app/service-user .
COPY .env .env

EXPOSE 8002

CMD ["./service-user"]