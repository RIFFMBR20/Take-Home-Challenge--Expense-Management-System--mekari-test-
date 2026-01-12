# Stage 1: Build
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go mod dan sum dulu untuk caching layer
COPY go.mod go.sum ./
RUN go mod download

# Copy all code
COPY . .

# Build app
RUN go build -o main .

# Stage 2: Run
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy binary from stage builder
COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]