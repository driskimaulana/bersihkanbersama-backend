# Start from a base Go image
FROM golang:1.16-alpine as builder

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum files to download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Start a new, lightweight image
FROM alpine:latest

# Set the working directory
WORKDIR /app

# Copy the built executable from the previous image
COPY --from=builder /app/main .

# Expose the port your Go application listens on (default is 8080)
EXPOSE 8080

# Run the Go application
CMD ["./main"]