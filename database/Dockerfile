FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY . .

RUN go work sync

COPY database/.env database/.env

WORKDIR /app/database

RUN go build -o database .

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/database/database .
COPY --from=builder /app/database/.env .

CMD ["./database"]