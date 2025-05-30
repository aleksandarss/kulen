# 1. Use a tiny builder image
FROM golang:1.24-alpine AS builder

# 2. Set working dir and copy go files
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# 3. Build the binary
RUN go build -o app ./cmd/server

# 4. Final image: tiny and secure
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/app .
# COPY .env .env

# 5. Expose port and run
EXPOSE 8080
CMD ["./app"]
