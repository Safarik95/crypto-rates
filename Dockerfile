FROM golang:1.24-alpine AS builder
RUN apk add --no-cache git
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o server ./cmd/service

FROM alpine:latest
RUN apk add --no-cache ca-certificates tzdata
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
WORKDIR /app
COPY --from=builder /app/server .
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/docs ./docs
RUN chown -R appuser:appgroup /app
USER appuser
EXPOSE 8080
CMD ["./server"]