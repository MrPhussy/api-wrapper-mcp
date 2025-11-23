FROM golang:alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o api_wrapper main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/api_wrapper .
COPY config ./config

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

COPY start.sh .
RUN chmod +x api_wrapper start.sh

# Default command
ENTRYPOINT ["./start.sh"]
