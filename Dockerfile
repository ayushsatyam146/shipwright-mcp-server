# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /workspace

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY main.go ./

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o shipwright-build-mcp-server .

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /workspace/shipwright-build-mcp-server .

# Run the binary
ENTRYPOINT ["./shipwright-build-mcp-server"] 