# Use the official Golang Alpine image as a builder.
FROM golang:1.18-alpine as builder

# Install git, required for fetching Go dependencies.
RUN apk add --no-cache git

# Set the Current Working Directory inside the container.
WORKDIR /app

# Copy go mod and sum files.
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed.
RUN go mod download

# Copy the source code into the container.
COPY . .

# Build the Go app as a static binary.
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o myapp .

# Final stage: Use Alpine for the production image.
FROM alpine:latest

# Install CA certificates for HTTPS connections.
RUN apk --no-cache add ca-certificates

# Copy the binary from the builder stage.
COPY --from=builder /app/myapp /myapp

# Command to run the executable.
CMD ["/myapp"]
