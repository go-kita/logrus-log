# ==============================================================================
# Makefile helper functions for tools
#

CANDIDATE_DEP_TOOLS := swagger mockgen gotests gsemver git-chglog github-release coscmd protoc-gen-go cfssl addlicense
DEP_TOOLS ?= golangci-lint go-junit-report golines go-mod-outdated goimports go-gitlint
ifeq ($(GOOS),darwin)
	DEP_TOOLS += gawk gxargs
endif
CANDIDATE_OTHER_TOOLS := depth go-callvis gothanks richgo rts
OTHER_TOOLS ?=

tools.install:
	@for tool in $(DEP_TOOLS) $(OTHER_TOOLS); do \
		echo "===========> Installing $$tool"; \
		$(MAKE) tools.verify.$$tool; \
	done

tools.verify.%:
	@if ! which $* &>/dev/null; then $(MAKE) install.$*; fi

.PHONY: install.swagger
install.swagger:
	@$(GO) get -u github.com/go-swagger/go-swagger/cmd/swagger

.PHONY: intsall.go-gitlint
install.go-gitlint:
	@$(GO) get -u github.com/llorllale/go-gitlint/cmd/go-gitlint

.PHONY: install.golangci-lint
install.golangci-lint:
ifeq ($(GOOS),darwin)
	@HOMEBREW_NO_AUTO_UPDATE=1 brew install -q golangci-lint
else
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | \
	sh -s -- -b $$($(GO) env GOPATH)/bin v1.41.1
endif

.PHONY: install.go-junit-report
install.go-junit-report:
	@$(GO) get -u github.com/jstemmer/go-junit-report

.PHONY: install.gawk # macOS need only
install.gawk:
ifeq ($(GOOS),darwin)
	@HOMEBREW_NO_AUTO_UPDATE=1 brew install -q gawk
else
	@echo "You need not install gawk, for your OS: $(GOOS)_$(GOARCH)"
endif

.PHONY: install.golines
install.golines:
	@$(GO) get -u github.com/segmentio/golines

.PHONY: install.go-mod-outdated
install.go-mod-outdated:
	@$(GO) get -u github.com/psampaz/go-mod-outdated

.PHONY: install.goimports
install.goimports:
	@$(GO) get -u golang.org/x/tools/cmd/goimports

.PHONY: install.gxargs
install.gxargs:
ifeq ($(GOOS),darwin)
	@HOMEBREW_NO_AUTO_UPDATE=1 brew install -q findutils
else
	@echo "You need not install findutils, for your OS: $(GOOS)_$(GOARCH)"
endif
