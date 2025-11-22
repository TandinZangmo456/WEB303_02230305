FROM golang:1.23-alpine AS builder
WORKDIR /app
ENV GOTOOLCHAIN=auto

# Copy proto files
COPY proto ./proto

# Copy gateway files
COPY api-gateway/go.mod api-gateway/go.sum ./
RUN go mod download

COPY api-gateway ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/server .
EXPOSE 8080
CMD ["/app/server"]