FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o ai-service ./cmd/grpc

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/ai-service .
COPY --from=builder /app/configs ./configs

EXPOSE 50051

CMD ["./ai-service"]