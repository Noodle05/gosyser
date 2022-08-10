FROM golang:1.19-alpine3.16 as builder

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . ./

RUN go build -v -o server /app/cmd/syslog

FROM alpine:3.16

ENV APP_CONFIG_FILE=/etc/syslog-server/application.yml
ENV LOG_LEVEL=info

RUN apk add --no-cache tzdata

COPY --from=builder /app/server /usr/sbin/syslog-server
COPY configs/application.yml "$APP_CONFIG_FILE"

EXPOSE 514/tcp
EXPOSE 514/udp

CMD ["/usr/sbin/syslog-server"]