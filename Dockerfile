# Start from the official Go image
FROM golang:1.24

WORKDIR /app

RUN apt-get update && \
    apt-get install -y poppler-utils && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# Copy go.mod and go.sum and download dependencies
COPY ./src/go.mod ./src/go.sum ./
RUN go mod download

# Copy the full project
COPY ./src .

# RUN go build -o /bin/building_embedding_generation ./cmd/building_embedding_generation.go
# RUN go build -o /bin/contact_embedding_generation ./cmd/contact_embedding_generation.go
RUN go build -o /bin/worker ./cmd/queue/worker.go
RUN go build -o /bin/logispro

# Default command (can be overridden in docker-compose)
CMD ["echo", "Use docker-compose to run the appropriate worker"]
