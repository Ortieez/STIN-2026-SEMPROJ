# Build stage
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Copy go mod and sum files from backend directory
COPY backend/go.mod backend/go.sum ./

# Download all dependencies
RUN go mod download

# Copy the backend source
COPY backend/ .

# Copy .env from root context
COPY .env ./.env

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the Pre-built binary
COPY --from=builder /app/main .
COPY --from=builder /app/.env .

# Expose port 8080
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
