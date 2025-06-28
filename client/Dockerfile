# Builder stage
FROM golang:1.22-alpine AS builder
WORKDIR /app

# Copy module files first
COPY go.mod go.sum ./
RUN go mod download

# Copy all source including protocols and benchmarks
COPY . .

# Build client
RUN CGO_ENABLED=0 GOOS=linux go build -o client ./client/main.go

# Runtime stage
FROM alpine:3.19
WORKDIR /app

# Copy binary
COPY --from=builder /app/client .
COPY --from=builder /app/protocols ./protocols
COPY --from=builder /app/benchmarks ./benchmarks

CMD ["./client"]