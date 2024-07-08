# Use the official Python slim image as a base
FROM python:3.9-slim

# Install Go and curl
RUN apt-get update && apt-get install -y \
    golang-go \
    curl \
    build-essential \
    libjpeg-dev \
    zlib1g-dev \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copy application files
COPY . .

# Install Python dependencies
RUN pip install pillow

# Install Go dependencies
RUN go get github.com/confluentinc/confluent-kafka-go/kafka
RUN go get github.com/gin-gonic/gin
RUN go get github.com/go-redis/redis/v8

# Build Go application
RUN go build -o server main.go

# Set environment variables for Kafka and Redis
ENV KAFKA_BROKER=localhost:9092
ENV REDIS_ADDR=redis:6379
ENV REDIS_PASS=""

# Expose the port the app runs on
EXPOSE 8080

# Run the Go server by default when the container starts
CMD ["/app/server"]
