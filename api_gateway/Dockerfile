FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY . .

RUN go work sync

COPY api_gateway/.env api_gateway/.env

WORKDIR /app/api_gateway

RUN go build -o server .

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/api_gateway/server .
COPY --from=builder /app/api_gateway/.env .

CMD ["./server"]
