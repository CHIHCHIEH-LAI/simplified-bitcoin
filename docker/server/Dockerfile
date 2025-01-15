# Use the Go base image
FROM golang:1.23

# Set the working directory inside the container
WORKDIR /app

# Copy all project files into the container
COPY . .

# Install dependencies
RUN go mod download

# Build the server binary
RUN go build -o server ./cmd/server/main.go

# Set environment variables (can be overridden at runtime)
ENV PORT=8080
ENV ADDRESS=127.0.0.1:8080
ENV BOOTSTRAP=""

# Expose the port the server will run on
EXPOSE $(PORT)

# Command to start the node with environment variables
CMD ["sh", "-c", "./server -port=${PORT} -address=${ADDRESS} -bootstrap=${BOOTSTRAP}"]