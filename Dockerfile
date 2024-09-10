FROM golang:1.13-alpine AS builder

RUN apk add --no-cache make gcc musl-dev linux-headers git

ADD . /tomochain
RUN cd /tomochain && make viction

FROM alpine:latest

WORKDIR /tomochain

COPY --from=builder /tomochain/build/bin/viction /usr/local/bin/viction

RUN chmod +x /usr/local/bin/viction

EXPOSE 8545
EXPOSE 30303

ENTRYPOINT ["/usr/local/bin/viction"]

CMD ["--help"]
