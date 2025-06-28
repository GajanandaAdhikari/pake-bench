# Builder stage
FROM golang:1.22-alpine AS builder
WORKDIR /app

# Copy module files first
COPY go.mod go.sum ./
RUN go mod download

# Copy all source including protocols and benchmarks
COPY . .

# Build server
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./server/main.go

# Runtime stage
FROM alpine:3.19
WORKDIR /app

# Copy binary and shared code
COPY --from=builder /app/server .
COPY --from=builder /app/protocols ./protocols
COPY --from=builder /app/benchmarks ./benchmarks

EXPOSE 8080
CMD ["./server"]