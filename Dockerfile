# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o janus-ssh .

# Runtime stage
FROM alpine:latest

# Install runtime dependencies (openssh-client is needed for SSH proxying)
RUN apk add --no-cache openssh-client ca-certificates

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/janus-ssh .

# Create directories for keys and data
RUN mkdir -p /app/keys /app/data

# Set proper permissions
RUN chmod 700 /app/keys

EXPOSE 2222

ENTRYPOINT ["./janus-ssh"]
