FROM golang:1.17 as builder

WORKDIR /sensorapi
COPY . /sensorapi

RUN export CGO_ENABLED=0 && go build -o sensorapi ./

FROM alpine:3.15
RUN apk update && apk add --no-cache bash
COPY --from=builder /sensorapi /sensorapi

VOLUME "/var/lib/server"

ENV PORT 8081

EXPOSE ${PORT}/tcp

WORKDIR /sensorapi
ENTRYPOINT ["./sensorapi", "start"]