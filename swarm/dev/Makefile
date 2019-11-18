.PHONY: build cluster test

default: build

build:
	go build -o bin/swarm    github.com/tomochain/tomochain/cmd/swarm
	go build -o bin/tomo     github.com/tomochain/tomochain/cmd/tomo
	go build -o bin/bootnode github.com/tomochain/tomochain/cmd/bootnode

cluster: build
	scripts/boot-cluster.sh

test:
	go test -v github.com/tomochain/tomochain/swarm/...
