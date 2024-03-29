BINARY_NAME = gojs
PACKAGE = github.com/ruiconti/gojs
PACKAGE_LEXER = $(PACKAGE)/lexer
PACKAGE_PARSER = $(PACKAGE)/parser
BUILD_FLAGS =-race -trimpath
ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))

build:
	@echo "Building the binary..."
	go build $(BUILD_FLAGS) -o $(BINARY_NAME) $(PACKAGE)

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
	
astviz:
	@go build -o bin/astviz $(PACKAGE)/tools
	@./bin/astviz $(ARGS)

# PHONY targets
.PHONY: build test clean astviz