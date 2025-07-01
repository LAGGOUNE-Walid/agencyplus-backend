# Start from the official Go image
FROM golang:1.24

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files and download deps
COPY ./src/go.mod ./src/go.sum ./
RUN go mod download

# Copy the rest of the app
COPY ./src .