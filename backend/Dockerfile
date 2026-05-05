# Build stage
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the Pre-built binary
COPY --from=builder /app/main .

# Expose port 8080
EXPOSE 8080

# Command to run the executable
# Using sh -c to ensure we see output and potentially handle signal forwarding
CMD ["./main"]
