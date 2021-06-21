.DEFAULT := all

.PHONY: all
all: format lint cover build

ROOT_PACKAGE=github.com/go-kita/log
VERSION_PACKAGE=

# ==============================================================================
# Includes

include scripts/make-rules/common.mk
include scripts/make-rules/golang.mk
include scripts/make-rules/tools.mk

# ==============================================================================
# Targets

## build
.PHONY: build
build:
	@$(MAKE) go.build.lib

## clean: Remove all files that are created by building.
.PHONY: clean
clean:
	@$(MAKE) go.clean

## lint: Check syntax and styling of go sources.
.PHONY: lint
lint: lint
	@$(MAKE) go.lint

## test: Run unit test.
.PHONY: test
test:
	@$(MAKE) go.test

## cover: Run unit test and get test coverage.
.PHONY: cover
cover:
	@$(MAKE) go.test.cover

## format: Gofmt (reformat) package sources (exclude vendor dir if existed).
.PHONY: format
format: tools.verify.golines tools.verify.goimports
	@echo "===========> Formating codes"
	@$(FIND) -type f -name '*.go' | $(XARGS) gofmt -s -w
	@$(FIND) -type f -name '*.go' | $(XARGS) goimports -w -local $(ROOT_PACKAGE)
	@$(FIND) -type f -name '*.go' | $(XARGS) golines -w --max-len=120 --reformat-tags --shorten-comments --ignore-generated .

## tools: install dependent tools.
.PHONY: tools
tools:
	@$(MAKE) tools.install
	@chmod +x $(ROOT_DIR)/githooks/commit-msg
	@cp -f $(ROOT_DIR)/githooks/commit-msg $(ROOT_DIR)/.git/hooks/commit-msg

## check-updates: Check outdated dependencies of the go projects.
.PHONY: check-updates
check-updates:
	@$(MAKE) go.updates
