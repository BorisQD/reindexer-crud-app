# Build
FROM golang:1.24 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /main ./cmd/main.go

# Runtime
FROM alpine:latest
WORKDIR /root/

COPY --from=builder /main .

CMD ["./main"]
