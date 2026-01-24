FROM golang:1.24 AS builder

WORKDIR /app

# Install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
# CGO_ENABLED=1 is often required for go-sqlite3, though we aim to use postgres specific usage might trigger it.
# We build a static binary if possible, but with CGO it's dynamic.
RUN CGO_ENABLED=1 go build -o main .

# Final stage
FROM debian:bookworm-slim

WORKDIR /app

# Install ca-certificates for external API calls (e.g. OpenAI)
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

# Copy binary and config
COPY --from=builder /app/main .
COPY --from=builder /app/config ./config
COPY --from=builder /app/.env.example ./.env.example

# Create storage directory
RUN mkdir -p storage/database storage/uploads

EXPOSE 8888

CMD ["./main"]
