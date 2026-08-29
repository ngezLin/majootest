# Stage 1: Build
FROM golang:alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build static binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/bin/api ./cmd

# Stage 2: Final Minimal Image
FROM alpine:3.20

WORKDIR /app

RUN apk --no-cache add ca-certificates tzdata curl

# Copy compiled binary from builder
COPY --from=builder /app/bin/api /app/api
COPY --from=builder /app/migrations /app/migrations
COPY --from=builder /app/docs /app/docs

EXPOSE 8080

CMD ["/app/api"]
