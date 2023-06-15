# Use the Go 1.14 base image
FROM golang:1.19.4

# Set the working directory inside the container
WORKDIR /cmd

# Copy the Go module files
COPY ./ ./

# Download the Go module dependencies
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go application
RUN go build -o main cmd/internal/main.go

# Set the entry point of the container
ENTRYPOINT ["./cmd/internal/fully_check.go"]

# Expose the port on which the application listens
EXPOSE 50051
