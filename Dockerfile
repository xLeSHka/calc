FROM golang:1.23.4-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o calc_service ./cmd/main

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/calc_service /app/calc_service

CMD /app/calc_service