FROM golang:1.25-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY src/ src/
RUN go build -o indervirdev ./src/

FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/indervirdev .

EXPOSE 22

CMD ["./indervirdev"]
