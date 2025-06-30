FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY . .

RUN go work sync

COPY message/.env message/.env

WORKDIR /app/message

RUN go build -o message .

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/message/message .
COPY --from=builder /app/message/.env .

CMD ["./message"]