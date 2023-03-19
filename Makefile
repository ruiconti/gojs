# Makefile for github.com/ruiconti/gojs

BINARY_NAME = gojs
PACKAGE = github.com/ruiconti/gojs
PACKAGE_LEXER = lexer
PACKAGE_PARSER = parser
BUILD_FLAGS = -race -trimpath

# Build target
build:
	@echo "Building the binary..."
	@go $BUILD_FLAGS build -o $(BINARY_NAME) $(PACKAGE)

# Test target
test:
	@echo "Running tests..."
	@go test $(PACKAGE)/$(PACKAGE_LEXER)
	@go test $(PACKAGE)/$(PACKAGE_PARSER)

# Clean target
clean:
	@echo "Cleaning up..."
	@go clean -modcache
	@rm -f $(BINARY_NAME)
	
# Tag targets
bump-pkg:
	@echo "Latest version:";
	@git tag --list "$(PACKAGE)/*" --sort=-version:refname | head -n 1; \
	read -p "Enter bump version for $(PACKAGE) (e.g. v1.0.0): " version; \
	git tag $(PACKAGE)/$$version; \
	git push origin $(PACKAGE)/$$version

bump-parser:
	@echo "Latest version:"; \
	@git tag --list "$(PACKAGE_PARSER)/*" --sort=-version:refname | head -n 1 || echo "No tags found"; \
	read -p "Enter bump version for $(PACKAGE_PARSER) (e.g. v1.0.0): " version; \
	git tag $(PACKAGE_PARSER)/$$version; \
	git push origin $(PACKAGE_PARSER)/$$version

bump-lexer:
	@echo "Latest version:"; \
	@git tag --list "$(PACKAGE_LEXER)/*" --sort=-version:refname | head -n 1 || echo "No tags found"; \
	read -p "Enter bump version $(PACKAGE_LEXER) (e.g. v1.0.0): " version; \
	git tag $(PACKAGE_LEXER)/$$version; \
	git push origin $(PACKAGE_LEXER)/$$version

# PHONY targets
.PHONY: build test clean