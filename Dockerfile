# Build stage
FROM golang:1.23.1-alpine AS builder

# Install build dependencies
RUN apk add --no-cache gcc musl-dev sqlite-dev

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o forum .

# Final stage
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates sqlite-libs

# Create app directory
WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/forum .

# Copy necessary files and directories
COPY --from=builder /app/web ./web
COPY --from=builder /app/database ./database
COPY --from=builder /app/forum.db ./forum.db

# Expose port
EXPOSE 7777

# Run the application
CMD ["./forum"]
