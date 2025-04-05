# Start with a minimal Go image for building the application
FROM golang:1.24 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go source code into the container
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go web application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server cmd/server/main.go

# Use a minimal base image for the final container
FROM alpine:latest

# Set the working directory in the new container
WORKDIR /root/

# Download a binary from GitHub releases
RUN apk add --no-cache curl && \
    curl -L -o migrate.linux-amd64.tar.gz https://github.com/golang-migrate/migrate/releases/download/v4.18.2/migrate.linux-amd64.tar.gz && \
    tar -xzf migrate.linux-amd64.tar.gz && \
    rm migrate.linux-amd64.tar.gz 

# Copy the compiled binary from the builder stage
COPY --from=builder /app/server .

# Expose the port the app runs on
EXPOSE 3000

# Run the application
CMD ["./server"]
