FROM golang:1.22-alpine AS builder

RUN apk --no-cache add ca-certificates tzdata git

WORKDIR /app

# Copy source first so go mod tidy can analyse imports
COPY . .

RUN go mod tidy && \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o server .

# ── minimal runtime image ─────────────────────────────────────
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

COPY --from=builder /app/server .

EXPOSE 8080

CMD ["./server"]
