# Use Go 1.24 as the builder stage
FROM golang:1.24 AS builder
WORKDIR /app

# Set Go module proxy
ENV GOPROXY="https://mirrors.aliyun.com/goproxy/"
ENV GO111MODULE=on
# Disable CGO for static compilation
ENV CGO_ENABLED=0

# Copy the source code
COPY . .

# Build the static binary
RUN go build -o quota_exporter -v ./cmd/main.go && chmod +x quota_exporter

# Use Alpine as the final runtime environment
FROM alpine:latest
WORKDIR /root/

# Install necessary libraries
RUN sed -i 's|dl-cdn.alpinelinux.org|mirrors.aliyun.com|g' /etc/apk/repositories \
    && apk add --no-cache ca-certificates

# Copy the built binary
COPY --from=builder /app/quota_exporter .

# Copy the configuration file
COPY config/config.yaml /root/config.yaml

# Start the application
CMD ["/root/quota_exporter"]
