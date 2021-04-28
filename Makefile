GO_VERSION=1.16
GOFILES = $(shell find . -type f -name '*.go' -not -path "./.git/*")
LDFLAGS = '-s -w -extldflags "-static" -X github.com/gimlet-io/gimlet-stack/version.Version='${VERSION}

DOCKER_RUN?=
_with-docker:
	$(eval DOCKER_RUN=docker run --rm -v $(shell pwd):/go/src/github.com/gimlet-io/gimlet-stack -w /go/src/github.com/gimlet-io/gimlet-stack golang:$(GO_VERSION))

.PHONY: all format test build dist fast-dist

all: test build

format:
	@gofmt -w ${GOFILES}

test:
	$(DOCKER_RUN) go test -race -timeout 60s $(shell go list ./... )

build:
	$(DOCKER_RUN) CGO_ENABLED=0 go build -ldflags $(LDFLAGS) -o build/gimlet github.com/gimlet-io/gimlet-stack/cmd

dist:
	mkdir -p bin
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags $(LDFLAGS) -a -installsuffix cgo -o bin/stack-linux-x86_64 github.com/gimlet-io/gimlet-stack/cmd
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags $(LDFLAGS) -a -installsuffix cgo -o bin/stack-darwin-x86_64 github.com/gimlet-io/gimlet-stack/cmd
	CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=6 go build -ldflags $(LDFLAGS) -a -installsuffix cgo -o bin/stack-linux-armhf github.com/gimlet-io/gimlet-stack/cmd
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags $(LDFLAGS) -a -installsuffix cgo -o bin/stack-linux-arm64 github.com/gimlet-io/gimlet-stack/cmd
	CGO_ENABLED=0 GOOS=windows go build -ldflags $(LDFLAGS) -a -installsuffix cgo -o bin/stack.exe github.com/gimlet-io/gimlet-stack/cmd

fast-dist:
	mkdir -p bin
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags $(LDFLAGS) -a -installsuffix cgo -o bin/stack-linux-x86_64 github.com/gimlet-io/gimlet-stack/cmd
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags $(LDFLAGS) -a -installsuffix cgo -o bin/stack-darwin-x86_64 github.com/gimlet-io/gimlet-stack/cmd
