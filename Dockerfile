# Start from the official Go image
FROM golang:1.24

WORKDIR /app

# Copy go.mod and go.sum and download dependencies
COPY ./src/go.mod ./src/go.sum ./
RUN go mod download

# Copy the full project
COPY ./src .

# Build each worker binary
RUN go build -o bin/sms_worker_prepare ./cmd/sms_worker_prepare.go
RUN go build -o bin/sms_worker_send ./cmd/sms_worker_send.go
RUN go build -o bin/logispro

# Default command (can be overridden in docker-compose)
CMD ["echo", "Use docker-compose to run the appropriate worker"]
