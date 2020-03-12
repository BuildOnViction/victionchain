.PHONY: chancoin chancoin-cross evm all test clean
.PHONY: chancoin-linux chancoin-linux-386 chancoin-linux-amd64 chancoin-linux-mips64 chancoin-linux-mips64le
.PHONY: chancoin-darwin chancoin-darwin-386 chancoin-darwin-amd64

GOBIN = $(shell pwd)/build/bin
GOFMT = gofmt
GO ?= 1.12
GO_PACKAGES = .
GO_FILES := $(shell find $(shell go list -f '{{.Dir}}' $(GO_PACKAGES)) -name \*.go)

GIT = git

chancoin:
	go run build/ci.go install ./cmd/chancoin
	@echo "Done building."
	@echo "Run \"$(GOBIN)/chancoin\" to launch chancoin."

gc:
	go run build/ci.go install ./cmd/gc
	@echo "Done building."
	@echo "Run \"$(GOBIN)/gc\" to launch gc."

bootnode:
	go run build/ci.go install ./cmd/bootnode
	@echo "Done building."
	@echo "Run \"$(GOBIN)/bootnode\" to launch a bootnode."

puppeth:
	go run build/ci.go install ./cmd/puppeth
	@echo "Done building."
	@echo "Run \"$(GOBIN)/puppeth\" to launch puppeth."

all:
	go run build/ci.go install

test: all
	go run build/ci.go test

clean:
	rm -fr build/_workspace/pkg/ $(GOBIN)/*

# Cross Compilation Targets (xgo)

chancoin-cross: chancoin-windows-amd64 chancoin-darwin-amd64 chancoin-linux
	@echo "Full cross compilation done:"
	@ls -ld $(GOBIN)/chancoin-*

chancoin-linux: chancoin-linux-386 chancoin-linux-amd64 chancoin-linux-mips64 chancoin-linux-mips64le
	@echo "Linux cross compilation done:"
	@ls -ld $(GOBIN)/chancoin-linux-*

chancoin-linux-386:
	go run build/ci.go xgo -- --go=$(GO) --targets=linux/386 -v ./cmd/chancoin
	@echo "Linux 386 cross compilation done:"
	@ls -ld $(GOBIN)/chancoin-linux-* | grep 386

chancoin-linux-amd64:
	go run build/ci.go xgo -- --go=$(GO) --targets=linux/amd64 -v ./cmd/chancoin
	@echo "Linux amd64 cross compilation done:"
	@ls -ld $(GOBIN)/chancoin-linux-* | grep amd64

chancoin-linux-mips:
	go run build/ci.go xgo -- --go=$(GO) --targets=linux/mips --ldflags '-extldflags "-static"' -v ./cmd/chancoin
	@echo "Linux MIPS cross compilation done:"
	@ls -ld $(GOBIN)/chancoin-linux-* | grep mips

chancoin-linux-mipsle:
	go run build/ci.go xgo -- --go=$(GO) --targets=linux/mipsle --ldflags '-extldflags "-static"' -v ./cmd/chancoin
	@echo "Linux MIPSle cross compilation done:"
	@ls -ld $(GOBIN)/chancoin-linux-* | grep mipsle

chancoin-linux-mips64:
	go run build/ci.go xgo -- --go=$(GO) --targets=linux/mips64 --ldflags '-extldflags "-static"' -v ./cmd/chancoin
	@echo "Linux MIPS64 cross compilation done:"
	@ls -ld $(GOBIN)/chancoin-linux-* | grep mips64

chancoin-linux-mips64le:
	go run build/ci.go xgo -- --go=$(GO) --targets=linux/mips64le --ldflags '-extldflags "-static"' -v ./cmd/chancoin
	@echo "Linux MIPS64le cross compilation done:"
	@ls -ld $(GOBIN)/chancoin-linux-* | grep mips64le

chancoin-darwin: chancoin-darwin-386 chancoin-darwin-amd64
	@echo "Darwin cross compilation done:"
	@ls -ld $(GOBIN)/chancoin-darwin-*

chancoin-darwin-386:
	go run build/ci.go xgo -- --go=$(GO) --targets=darwin/386 -v ./cmd/chancoin
	@echo "Darwin 386 cross compilation done:"
	@ls -ld $(GOBIN)/chancoin-darwin-* | grep 386

chancoin-darwin-amd64:
	go run build/ci.go xgo -- --go=$(GO) --targets=darwin/amd64 -v ./cmd/chancoin
	@echo "Darwin amd64 cross compilation done:"
	@ls -ld $(GOBIN)/chancoin-darwin-* | grep amd64

chancoin-windows-amd64:
	go run build/ci.go xgo -- --go=$(GO) --targets=windows/amd64 -v ./cmd/chancoin
	@echo "Windows amd64 cross compilation done:"
	@ls -ld $(GOBIN)/chancoin-windows-* | grep amd64
gofmt:
	$(GOFMT) -s -w $(GO_FILES)
	$(GIT) checkout vendor
