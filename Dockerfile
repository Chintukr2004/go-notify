# Build Stage
FROM golang:alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /api ./cmd/api
RUN CGO_ENABLED=0 GOOS=linux go build -o /worker ./cmd/worker

# Final Stage for API Service
FROM alpine:latest AS api
WORKDIR /root/
COPY --from=builder /api .
EXPOSE 8080
CMD ["./api"]

# Final Stage for Worker Service
FROM alpine:latest AS worker
WORKDIR /root/
COPY --from=builder /worker .
CMD ["./worker"]