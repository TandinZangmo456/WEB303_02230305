FROM golang:1.23.2-alpine AS builder
WORKDIR /app
COPY time-service/go.mod ./
COPY time-service/go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./time-service/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/server .
EXPOSE 50052
CMD ["/app/server"]