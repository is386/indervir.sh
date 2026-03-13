FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY src/ src/
RUN go build -o indervirsh ./src/

FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/indervirsh .

EXPOSE 22

CMD ["./indervirsh"]
