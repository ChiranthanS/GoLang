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
COPY main.go .
COPY ai_model.py .
COPY images /app/images

# Install Python dependencies
RUN pip install pillow


# Build Go application
RUN go build -o server main.go

# Expose the port the app runs on
EXPOSE 8080

# Run the Go server by default when the container starts
CMD ["/app/server"]
