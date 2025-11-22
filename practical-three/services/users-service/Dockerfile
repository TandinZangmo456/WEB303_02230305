FROM golang:1.23-alpine AS builder
WORKDIR /app
ENV GOTOOLCHAIN=auto

# Copy proto files
COPY proto ./proto

# Copy service files
COPY services/users-service/go.mod services/users-service/go.sum ./
RUN go mod download

COPY services/users-service ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/server .
EXPOSE 50051
CMD ["/app/server"]