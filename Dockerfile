# Stage 1: Build with CGO
FROM golang:1.24.2 as builder

WORKDIR /app

ENV CGO_ENABLED=1 \
    GOOS=linux \
    GOARCH=amd64

RUN apt-get update && apt-get install -y gcc musl-dev libsqlite3-dev

COPY go.mod ./
RUN go mod download

COPY . .
COPY .env.example .env
RUN go build -o app ./cmd/app/main.go  # Adjust path to your actual main.go

# Stage 2: Runtime with pg_dump, mysqldump, sqlite3
FROM ubuntu:22.04

RUN apt-get update && apt-get install -y \
    sqlite3 \
    default-mysql-client \
    postgresql-client \
    && apt-get clean && rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/app /usr/local/bin/app
WORKDIR /app

ENTRYPOINT ["app"]
