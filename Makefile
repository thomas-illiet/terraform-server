SHELL := bash
NAME := terrapi
IMPORT := github.com/thomas-illiet/$(NAME)
BIN := bin
DIST := dist
UNAME := $(shell uname -s)

GOBUILD ?= CGO_ENABLED=0 go build
PACKAGES ?= $(shell go list ./...)
SOURCES ?= $(shell find . -name "*.go" -type f -not -path ./.devenv/\* -not -path ./.direnv/\*)
GENERATE ?= $(PACKAGES)
TAGS ?= netgo

.PHONY: generate
generate:
	go generate $(GENERATE)

.PHONY: vet
vet:
	go vet $(PACKAGES)

.PHONY: staticcheck
staticcheck: $(STATICCHECK)
	$(STATICCHECK) -tags '$(TAGS)' $(PACKAGES)

.PHONY: lint
lint: $(REVIVE)
	for PKG in $(PACKAGES); do $(REVIVE) -config revive.toml -set_exit_status $$PKG || exit 1; done;

.PHONY: build
build: $(BIN)/$(NAME)

$(BIN)/$(NAME): $(SOURCES)
	$(GOBUILD) -v -tags '$(TAGS)' -ldflags '$(LDFLAGS)' -o $@ ./cmd/$(NAME)

$(BIN)/$(NAME)-debug: $(SOURCES)
	$(GOBUILD) -v -tags '$(TAGS)' -ldflags '$(LDFLAGS)' -gcflags '$(GCFLAGS)' -o $@ ./cmd/$(NAME)

.PHONY: test
test:
	go test -coverprofile coverage.out $(PACKAGES)