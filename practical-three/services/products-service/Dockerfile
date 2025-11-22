FROM golang:1.23-alpine AS builder
WORKDIR /app
ENV GOTOOLCHAIN=auto

# Copy proto files
COPY proto ./proto

# Copy service files
COPY services/products-service/go.mod services/products-service/go.sum ./
RUN go mod download

COPY services/products-service ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/server .
EXPOSE 50052
CMD ["/app/server"]