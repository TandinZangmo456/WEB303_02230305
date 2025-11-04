FROM golang:1.23.2-alpine AS builder
WORKDIR /app
COPY greeter-service/go.mod ./
COPY greeter-service/go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./greeter-service/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/server .
EXPOSE 50051
CMD ["/app/server"]