# Use the official Golang image as the base image
FROM golang:1.22.3 AS builder

# Set the working directory inside the container
WORKDIR /

# Copy the source code into the container
COPY . .

# Download the dependencies
RUN go mod download

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

# Create a new lightweight image with only the binary
FROM scratch

# Copy the binary from the builder stage into the new image
COPY --from=builder /app /app

# Set the entry point for the container
ENTRYPOINT ["/app"]
