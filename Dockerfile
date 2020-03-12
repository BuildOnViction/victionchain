FROM golang:1.12-alpine as builder

RUN apk add --no-cache make gcc musl-dev linux-headers git

ADD . /chancoin
RUN cd /chancoin && make chancoin

FROM alpine:latest

WORKDIR /chancoin

COPY --from=builder /chancoin/build/bin/chancoin /usr/local/bin/chancoin

RUN chmod +x /usr/local/bin/chancoin

EXPOSE 8545
EXPOSE 30303

ENTRYPOINT ["/usr/local/bin/chancoin"]

CMD ["--help"]
