# base go image
FROM golang:1.18-alpine as builder

# Update certificate for Go's dependency registry
RUN apk update && apk add ca-certificates && apk add openssl
RUN openssl s_client -showcerts -connect proxy.golang.org:443 | openssl x509 -outform PEM > /usr/local/share/ca-certificates/proxy.golang.crt
RUN update-ca-certificates

RUN mkdir /app
COPY . /app
WORKDIR /app

RUN CGO_ENABLED=0 go build -o brokerApp ./cmd/api

RUN chmod +x /app/brokerApp

# build a tiny docker image

FROM alpine:latest

RUN mkdir /app

COPY --from=builder /app/brokerApp /app

CMD [ "/app/brokerApp" ]