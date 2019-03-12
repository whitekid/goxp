.PHONY: clean test get tidy

get:
	go get -u
	@$(MAKE) tidy

tidy:
	go mod tidy
