.PHONY: build cluster test

default: build

build:
	go build -o bin/swarm    github.com/chancoin-core/chancoin-gold/cmd/swarm
	go build -o bin/chancoin     github.com/chancoin-core/chancoin-gold/cmd/chancoin
	go build -o bin/bootnode github.com/chancoin-core/chancoin-gold/cmd/bootnode

cluster: build
	scripts/boot-cluster.sh

test:
	go test -v github.com/chancoin-core/chancoin-gold/swarm/...
