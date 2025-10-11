# Development stage
FROM golang:1.21-alpine AS development

WORKDIR /app

# Install air for hot reload
RUN go install github.com/cosmtrek/air@latest

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

EXPOSE 8080

# Use air for hot reload in development
CMD ["air", "-c", ".air.toml"]

# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server

# Production stage
FROM alpine:latest AS production

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/migrations ./migrations

EXPOSE 8080

CMD ["./main"]
