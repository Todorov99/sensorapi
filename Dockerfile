FROM golang:1.17 as builder

WORKDIR /httpServer

COPY . /httpServer

RUN go build -o server .

FROM alpine:3.11
RUN apk update && apk add bash
COPY --from=builder /httpServer /httpServer

VOLUME "/var/lib/server"

EXPOSE 8081/tcp
ENTRYPOINT ["./httpServer/server"]
CMD [ "./httpServer/server" ]