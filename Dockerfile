# Stage 1: Build the application
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Enable CGO_ENABLED=0 for static binary
ENV CGO_ENABLED=0 GOOS=linux

# Copy go mod and sum files
COPY go.mod ./
# COPY go.sum ./ # Uncomment if you have dependencies

# Download all dependencies.
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN go build -o main .

# Stage 2: Run the application
FROM alpine:latest

WORKDIR /app/

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main .

# Expose port 8080 (optional but good documentation)
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
