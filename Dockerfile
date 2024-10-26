FROM golang:1.18-alpine AS builder

RUN apk add --no-cache make gcc musl-dev linux-headers git

ADD . /tomochain
RUN cd /tomochain && make tomo

FROM alpine:latest

WORKDIR /tomochain

COPY --from=builder /tomochain/build/bin/tomo /usr/local/bin/tomo

RUN chmod +x /usr/local/bin/tomo

EXPOSE 8545
EXPOSE 30303

ENTRYPOINT ["/usr/local/bin/tomo"]

CMD ["--help"]
