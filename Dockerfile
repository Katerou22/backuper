# Stage 1: Build the Go binary
FROM golang:1.21-alpine AS builder

# Install git (needed for Go modules), and set up working directory
RUN apk add --no-cache git
WORKDIR /app

# Copy go mod files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source and build
COPY . .
RUN go build -o app ./cmd/app/main.go  # Adjust path if needed

# Stage 2: Final image with dump tools and app
FROM alpine:latest

# Install latest pg_dump and mysqldump
RUN apk add --no-cache postgresql-client mysql-client

# Copy built binary from builder stage
COPY --from=builder /app/app /usr/local/bin/app



# Set working directory and run the app
WORKDIR /app
ENTRYPOINT ["app"]
