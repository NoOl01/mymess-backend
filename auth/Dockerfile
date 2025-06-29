FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY . .

RUN go work sync

COPY auth/.env auth/.env

WORKDIR /app/auth

RUN go build -o auth .

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/auth/auth .
COPY --from=builder /app/auth/.env .

CMD ["./auth"]