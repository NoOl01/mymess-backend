FROM golang:1.24 AS builder

WORKDIR /app

COPY . .

RUN go work sync

ARG BUILD_SERVICE

WORKDIR /app/${BUILD_SERVICE}

RUN CGO_ENABLED=1 go build -o ${BUILD_SERVICE} .

FROM debian:stable-slim

WORKDIR /root/

ARG BUILD_SERVICE
ENV SERVICE_BINARY=./${BUILD_SERVICE}

COPY --from=builder /app/${BUILD_SERVICE}/${BUILD_SERVICE} .
COPY --from=builder /app/.env .

RUN chmod +x ${BUILD_SERVICE}

CMD sh -c "$SERVICE_BINARY"