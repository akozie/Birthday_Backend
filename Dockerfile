# Stage 1: Build
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main cmd/main.go

# Stage 2: Run
FROM alpine:latest
WORKDIR /app
# Only copy the binary from the builder
COPY --from=builder /app/main .

# Expose port 8080
EXPOSE 8080
CMD ["./main"]