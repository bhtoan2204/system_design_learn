GO ?= go
GOCACHE ?= $(CURDIR)/.gocache

.PHONY: build run test fmt swagger clean

build:
	GOCACHE=$(GOCACHE) $(GO) build ./...

run:
	GOCACHE=$(GOCACHE) $(GO) run .

test:
	GOCACHE=$(GOCACHE) $(GO) test ./...

fmt:
	@gofmt -w $(shell find . -name '*.go' -not -path './.git/*')

swagger:
	@echo "Swagger spec available at http://localhost:8080/swagger.yaml"
	@echo "Local file: docs/openapi.yaml"

clean:
	rm -rf $(GOCACHE)
