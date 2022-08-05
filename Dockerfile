FROM golang:1.19-alpine3.16 as builder

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . ./

RUN go build -v -o server /app/cmd/syslog

FROM alpine:3.16

ENV APP_CONFIG_FILE=/etc/syslog-server/application.yml

RUN apk add --no-cache tzdata

COPY --from=builder /app/server /usr/sbin/syslog-server
COPY configs/application.yml "$APP_CONFIG_FILE"

CMD ["/usr/sbin/syslog-server"]