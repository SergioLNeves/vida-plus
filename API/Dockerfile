# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./cmd/api

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy the binary and environment file from builder
COPY --from=builder /app/main .
COPY --from=builder /app/.env.example .env

EXPOSE 8080

# Run the application
CMD ["./main"]