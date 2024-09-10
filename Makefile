.PHONY: tomo vic-cross evm all test clean
.PHONY: vic-linux vic-linux-386 vic-linux-amd64 vic-linux-mips64 vic-linux-mips64le
.PHONY: vic-darwin vic-darwin-386 vic-darwin-amd64

GOBIN = $(shell pwd)/build/bin
GOFMT = gofmt
GO ?= 1.13.15
GO_PACKAGES = .
GO_FILES := $(shell find $(shell go list -f '{{.Dir}}' $(GO_PACKAGES)) -name \*.go)

GIT = git

viction:
	go run build/ci.go install ./cmd/viction
	@echo "Done building."
	@echo "Run \"$(GOBIN)/viction\" to launch Viction."

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

vic-cross: vic-windows-amd64 vic-darwin-amd64 vic-linux
	@echo "Full cross compilation done:"
	@ls -ld $(GOBIN)/vic-*

vic-linux: vic-linux-386 vic-linux-amd64 vic-linux-mips64 vic-linux-mips64le
	@echo "Linux cross compilation done:"
	@ls -ld $(GOBIN)/vic-linux-*

vic-linux-386:
	go run build/ci.go xgo -- --go=$(GO) --targets=linux/386 -v ./cmd/viction
	@echo "Linux 386 cross compilation done:"
	@ls -ld $(GOBIN)/vic-linux-* | grep 386

vic-linux-amd64:
	go run build/ci.go xgo -- --go=$(GO) --targets=linux/amd64 -v ./cmd/viction
	@echo "Linux amd64 cross compilation done:"
	@ls -ld $(GOBIN)/vic-linux-* | grep amd64

vic-linux-mips:
	go run build/ci.go xgo -- --go=$(GO) --targets=linux/mips --ldflags '-extldflags "-static"' -v ./cmd/viction
	@echo "Linux MIPS cross compilation done:"
	@ls -ld $(GOBIN)/vic-linux-* | grep mips

vic-linux-mipsle:
	go run build/ci.go xgo -- --go=$(GO) --targets=linux/mipsle --ldflags '-extldflags "-static"' -v ./cmd/viction
	@echo "Linux MIPSle cross compilation done:"
	@ls -ld $(GOBIN)/vic-linux-* | grep mipsle

vic-linux-mips64:
	go run build/ci.go xgo -- --go=$(GO) --targets=linux/mips64 --ldflags '-extldflags "-static"' -v ./cmd/viction
	@echo "Linux MIPS64 cross compilation done:"
	@ls -ld $(GOBIN)/vic-linux-* | grep mips64

vic-linux-mips64le:
	go run build/ci.go xgo -- --go=$(GO) --targets=linux/mips64le --ldflags '-extldflags "-static"' -v ./cmd/viction
	@echo "Linux MIPS64le cross compilation done:"
	@ls -ld $(GOBIN)/vic-linux-* | grep mips64le

vic-darwin: vic-darwin-386 vic-darwin-amd64
	@echo "Darwin cross compilation done:"
	@ls -ld $(GOBIN)/vic-darwin-*

vic-darwin-386:
	go run build/ci.go xgo -- --go=$(GO) --targets=darwin/386 -v ./cmd/viction
	@echo "Darwin 386 cross compilation done:"
	@ls -ld $(GOBIN)/vic-darwin-* | grep 386

vic-darwin-amd64:
	go run build/ci.go xgo -- --go=$(GO) --targets=darwin/amd64 -v ./cmd/viction
	@echo "Darwin amd64 cross compilation done:"
	@ls -ld $(GOBIN)/vic-darwin-* | grep amd64

vic-windows-amd64:
	go run build/ci.go xgo -- --go=$(GO) -buildmode=mode -x --targets=windows/amd64 -v ./cmd/viction
	@echo "Windows amd64 cross compilation done:"
	@ls -ld $(GOBIN)/vic-windows-* | grep amd64
gofmt:
	$(GOFMT) -s -w $(GO_FILES)
	$(GIT) checkout vendor
