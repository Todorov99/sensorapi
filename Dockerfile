FROM golang:1.17 as builder

WORKDIR /server
COPY . /server

RUN export CGO_ENABLED=0 && go build -o server ./

FROM alpine:3.15
RUN apk update && apk add --no-cache bash
COPY --from=builder /server /server

VOLUME "/var/lib/server"

EXPOSE 8081/tcp
WORKDIR /server
ENTRYPOINT ["./server", "start", "-p", "8081"]