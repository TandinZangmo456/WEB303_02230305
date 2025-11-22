FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /food-catalog-service .

FROM alpine:latest
RUN apk --no-cache add curl
WORKDIR /
COPY --from=builder /food-catalog-service /food-catalog-service
EXPOSE 8080
CMD ["/food-catalog-service"]