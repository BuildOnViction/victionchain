.PHONY: tomo tomo-cross evm all test clean
.PHONY: tomo-linux tomo-linux-386 tomo-linux-amd64 tomo-linux-mips64 tomo-linux-mips64le
.PHONY: tomo-darwin tomo-darwin-386 tomo-darwin-amd64

GOBIN = $(shell pwd)/build/bin
GOFMT = gofmt
GO ?= 1.13.1
GO_PACKAGES = .
GO_FILES := $(shell find $(shell go list -f '{{.Dir}}' $(GO_PACKAGES)) -name \*.go)

GIT = git

tomo:
	go run build/ci.go install ./cmd/tomo
	@echo "Done building."
	@echo "Run \"$(GOBIN)/tomo\" to launch tomo."

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

tomo-cross: tomo-windows-amd64 tomo-darwin-amd64 tomo-linux
	@echo "Full cross compilation done:"
	@ls -ld $(GOBIN)/tomo-*

tomo-linux: tomo-linux-386 tomo-linux-amd64 tomo-linux-mips64 tomo-linux-mips64le
	@echo "Linux cross compilation done:"
	@ls -ld $(GOBIN)/tomo-linux-*

tomo-linux-386:
	go run build/ci.go xgo -- --go=$(GO) --targets=linux/386 -v ./cmd/tomo
	@echo "Linux 386 cross compilation done:"
	@ls -ld $(GOBIN)/tomo-linux-* | grep 386

tomo-linux-amd64:
	go run build/ci.go xgo -- --go=$(GO) --targets=linux/amd64 -v ./cmd/tomo
	@echo "Linux amd64 cross compilation done:"
	@ls -ld $(GOBIN)/tomo-linux-* | grep amd64

tomo-linux-mips:
	go run build/ci.go xgo -- --go=$(GO) --targets=linux/mips --ldflags '-extldflags "-static"' -v ./cmd/tomo
	@echo "Linux MIPS cross compilation done:"
	@ls -ld $(GOBIN)/tomo-linux-* | grep mips

tomo-linux-mipsle:
	go run build/ci.go xgo -- --go=$(GO) --targets=linux/mipsle --ldflags '-extldflags "-static"' -v ./cmd/tomo
	@echo "Linux MIPSle cross compilation done:"
	@ls -ld $(GOBIN)/tomo-linux-* | grep mipsle

tomo-linux-mips64:
	go run build/ci.go xgo -- --go=$(GO) --targets=linux/mips64 --ldflags '-extldflags "-static"' -v ./cmd/tomo
	@echo "Linux MIPS64 cross compilation done:"
	@ls -ld $(GOBIN)/tomo-linux-* | grep mips64

tomo-linux-mips64le:
	go run build/ci.go xgo -- --go=$(GO) --targets=linux/mips64le --ldflags '-extldflags "-static"' -v ./cmd/tomo
	@echo "Linux MIPS64le cross compilation done:"
	@ls -ld $(GOBIN)/tomo-linux-* | grep mips64le

tomo-darwin: tomo-darwin-386 tomo-darwin-amd64
	@echo "Darwin cross compilation done:"
	@ls -ld $(GOBIN)/tomo-darwin-*

tomo-darwin-386:
	go run build/ci.go xgo -- --go=$(GO) --targets=darwin/386 -v ./cmd/tomo
	@echo "Darwin 386 cross compilation done:"
	@ls -ld $(GOBIN)/tomo-darwin-* | grep 386

tomo-darwin-amd64:
	go run build/ci.go xgo -- --go=$(GO) --targets=darwin/amd64 -v ./cmd/tomo
	@echo "Darwin amd64 cross compilation done:"
	@ls -ld $(GOBIN)/tomo-darwin-* | grep amd64

tomo-windows-amd64:
	go run build/ci.go xgo -- --go=$(GO) -buildmode=mode -x --targets=windows/amd64 -v ./cmd/tomo
	@echo "Windows amd64 cross compilation done:"
	@ls -ld $(GOBIN)/tomo-windows-* | grep amd64
gofmt:
	$(GOFMT) -s -w $(GO_FILES)
	$(GIT) checkout vendor
