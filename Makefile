# Makefile for github.com/ruiconti/gojs

BINARY_NAME = gojs
WORKSPACE = github.com/ruiconti/gojs
BUILD_FLAGS = -race -trimpath

# Build target
build:
	@echo "Building the binary..."
	@go $BUILD_FLAGS build -o $(BINARY_NAME) $(WORKSPACE)

# Test target
test:
	@echo "Running tests..."
	@go test -v ./...

# Clean target
clean:
	@echo "Cleaning up..."
	@rm -f $(BINARY_NAME)
	
# Tag targets
bump-pkg:
	git tag --list "$(WORKSPACE)/*" --sort=-version:refname | head -n 1; \
	read -p "Enter workspace version (e.g. v1.0.0): " version; \
	git tag $(WORKSPACE)/$$version; \
	git push origin $(WORKSPACE)/$$version

bump-parser:
	git tag --list "$(PACKAGE_PARSER)/*" --sort=-version:refname | head -n 1 || echo "No tags found"; \
	@read -p "Enter parser version (e.g. v1.0.0): " version; \
	git tag $(PACKAGE_PARSER)/$$version; \
	git push origin $(PACKAGE_PARSER)/$$version

bump-lexer:
	git tag --list "$(PACKAGE_LEXER)/*" --sort=-version:refname | head -n 1 || echo "No tags found"; \
	@read -p "Enter lexer version (e.g. v1.0.0): " version; \
	git tag $(PACKAGE_LEXER)/$$version; \
	git push origin $(PACKAGE_LEXER)/$$version

# Clean target
clean:
	@echo "Cleaning up..."
	@rm -f $(BINARY_NAME)

# PHONY targets
.PHONY: build test clean