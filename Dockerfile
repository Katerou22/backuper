# Stage 1: Build with CGO
FROM golang:1.24.2 as builder

WORKDIR /app

ENV CGO_ENABLED=1 \
    GOOS=linux \
    GOARCH=amd64

RUN apt-get update && apt-get install -y gcc musl-dev libsqlite3-dev

COPY go.mod ./
RUN go get \
    gorm.io/gorm \
    gorm.io/driver/sqlite \
    github.com/go-telegram/bot \
    github.com/robfig/cron/v3 \
    github.com/joho/godotenv
RUN go mod tidy
RUN go mod download


COPY . .
RUN go build -o app ./cmd/app/main.go  # Adjust path to your actual main.go

#Test
# Stage 2: Runtime with pg_dump, mysqldump, sqlite3
FROM ubuntu:latest

# Replace mirror to avoid archive.ubuntu.com errors
RUN sed -i 's|http://archive.ubuntu.com|http://mirror.ubuntu.com|g' /etc/apt/sources.list && \
    apt-get update && \
    apt-get install -y \
        sqlite3 \
        default-mysql-client \
        postgresql-client && \
    apt-get clean && rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/app /usr/local/bin/app
WORKDIR /app

ENTRYPOINT ["app"]
