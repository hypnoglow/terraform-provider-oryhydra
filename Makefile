.PHONY: all
all: deps build

.PHONY: deps
deps:
	@go mod download
	@go mod vendor
	@go mod tidy

.PHONY: build
build:
	@go build -o ./bin/terraform-provider-oryhydra

.PHONY: prepare-examples
prepare-examples:
	@ln -s $(shell pwd)/bin/terraform-provider-oryhydra ./examples/oryhydra_oauth2_client/terraform-provider-oryhydra

.PHONY: lint
lint:
	@golangci-lint run
	@golangci-lint fmt
