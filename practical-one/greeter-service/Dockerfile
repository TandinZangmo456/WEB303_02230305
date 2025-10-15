FROM golang:1.23.2-alpine AS builder
WORKDIR /app

# Copy go mod files
COPY greeter-service/go.mod greeter-service/go.sum ./greeter-service/

# Copy proto generated files
COPY proto/ ./proto/

# Download dependencies
WORKDIR /app/greeter-service
RUN go mod download

# Copy the service code
WORKDIR /app
COPY greeter-service/ ./greeter-service/

# Build
WORKDIR /app/greeter-service
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/server .
EXPOSE 50051
CMD ["/app/server"]