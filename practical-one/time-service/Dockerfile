FROM golang:1.23.2-alpine AS builder
WORKDIR /app

# Copy go mod files
COPY time-service/go.mod time-service/go.sum ./time-service/

# Copy proto generated files
COPY proto/ ./proto/

# Download dependencies
WORKDIR /app/time-service
RUN go mod download

# Copy the service code
WORKDIR /app
COPY time-service/ ./time-service/

# Build
WORKDIR /app/time-service
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/server .
EXPOSE 50052
CMD ["/app/server"]