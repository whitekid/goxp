TARGET=bin/goxp
SRC=$(shell find . -type f -name '*.go' -not -path "./vendor/*" -not -path "*_test.go")

GIT_COMMIT ?= $(shell git rev-parse HEAD)
GIT_SHA ?= $(shell git rev-parse --short HEAD)
GIT_BRANCH ?= $(shell git rev-parse --abbrev-ref HEAD)
GIT_TAG ?= $(shell git describe --tags --always)
GIT_DIRTY ?= $(shell test -n "`git status --porcelain`" && echo "dirty" || echo "clean")
VER_BUILD_TIME ?= $(shell date +%Y-%m-%dT%H:%M:%S%z)

LDFLAGS = -s -w
LDFLAGS += -X main.GitCommit=${GIT_COMMIT}
LDFLAGS += -X main.GitSHA=${GIT_SHA}
LDFLAGS += -X main.GitBranch=${GIT_BRANCH}
LDFLAGS += -X main.GitTag=${GIT_TAG}
LDFLAGS += -X main.GitDirty=${GIT_DIRTY}
LDFLAGS += -X main.BuildTime=${VER_BUILD_TIME}

GO?=go
BUILD_FLAGS?=-v -ldflags="${LDFLAGS}"

.PHONY: clean install test get tidy

all: build
build: $(TARGET)

$(TARGET): $(SRC)
	${GO} build -o bin/ ${BUILD_FLAGS} ./cmd/...

install:
	@go install -v ${BUILD_FLAGS} ./...

clean:
	rm -rf bin/

test:
	@go test -v ./...

dep:
	@rm -f go.mod go.sum
	@go mod init github.com/whitekid/goxp

	@$(MAKE) tidy

	@go mod edit -retract="[v0.0.1,v0.0.10]"
	@$(MAKE) tidy

tidy:
	@go mod tidy -v
