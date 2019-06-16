GO=go
NAME=cacheapi
TARGET=github.com/johejo/cacheapi/cmd/cacheapi

.PHONY: default
default: binary

.PHONY: clean
clean:
	rm -rf out/

.PHONY: binary
binary:
	$(GO) build -o out/$(NAME) $(TARGET)

.PHONY: test
test:
	$(GO) test -v ./...
