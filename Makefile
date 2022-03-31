TARGET=bin/goxp
SRC=$(shell find . -type f -name '*.go' -not -path "./vendor/*" -not -path "*_test.go")

GO?=go
BUILD_FLAGS?=-v

.PHONY: clean test get tidy

all: build
build: $(TARGET)

$(TARGET): $(SRC)
	${GO} build -o bin/ ${BUILD_FLAGS} ./cmd/...

test:
	@go test -v ./...

dep:
	@rm -f go.mod go.sum
	@go mod init github.com/whitekid/goxp

	@$(MAKE) tidy

tidy:
	@go mod tidy
