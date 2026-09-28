GO ?= go
TOOLS_MOD := -modfile=go.tools.mod
PACKAGES ?= $(shell $(GO) list ./...)
GOFILES := $(shell find . -name "*.go" -type f)

.PHONY: fmt
fmt:
	$(GO) tool $(TOOLS_MOD) golangci-lint fmt

.PHONY: fmt-check
fmt-check:
	@diff=$$($(GO) tool $(TOOLS_MOD) golangci-lint fmt --diff) || exit $$?; \
	if [ -n "$$diff" ]; then \
		echo "Please run 'make fmt' and commit the result:"; \
		echo "$${diff}"; \
		exit 1; \
	fi;

test: fmt-check
	@$(GO) test -v -cover -coverprofile coverage.txt ./... && echo "\n==>\033[32m Ok\033[m\n" || exit 1

vet:
	$(GO) vet ./...

clean:
	$(GO) clean -modcache -cache -i
	find . -name "coverage.txt" -delete

.PHONY: lint
lint: ## Run golangci-lint
	$(GO) tool $(TOOLS_MOD) golangci-lint run

.PHONY: install-tools fmt lint
install-tools: ## Download pinned Go tools
	$(GO) mod download $(TOOLS_MOD)
