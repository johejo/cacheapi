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

.PHONY: min-binary
min-binary:
	$(GO) build -ldflags '-s -w' -o out/$(NAME) $(TARGET)
	upx out/$(NAME)

.PHONY: docker
docker:
	docker build -t cacheapi .

.PHONY: test
test:
	$(GO) test -v ./...
