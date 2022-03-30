.PHONY: clean test get tidy

test:
	@go test -v ./...

dep:
	@rm -f go.mod go.sum
	@go mod init github.com/whitekid/goex

	@$(MAKE) tidy

tidy:
	@go mod tidy
