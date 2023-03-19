BINARY_NAME = gojs
PACKAGE = github.com/ruiconti/gojs
PACKAGE_LEXER = $(PACKAGE)/lexer
PACKAGE_PARSER = $(PACKAGE)/parser
BUILD_FLAGS = -race -trimpath

build:
	@echo "Building the binary..."
	@go $BUILD_FLAGS build -o $(BINARY_NAME) $(PACKAGE)

test:
	@echo "Running tests..."
	@go test $(PACKAGE_LEXER)
	@go test $(PACKAGE_PARSER)

clean:
	@echo "Cleaning up..."
	@go clean -modcache
	@rm -f $(BINARY_NAME)
	
bump:
	@echo "Latest version:";
	@git tag --list "$(PACKAGE)/*" --sort=-version:refname | head -n 1; \
	read -p "Enter bump version for $(PACKAGE) (e.g. v1.0.0): " version; \
	git tag $(PACKAGE)/$$version; \
	git push origin $(PACKAGE)/$$version

# PHONY targets
.PHONY: build test clean