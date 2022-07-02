.DEFAULT_TARGET=help
.PHONY: all
all: lint test

# COMMANDS
## run: runs the application locally
.PHONY: run
run:
	$(call blue, "# Running App...")
	@go run main.go

## test: run test suites
.PHONY: test
test:
	@go test -race ./... || (echo "go test failed $$?"; exit 1)

## lint: run golangci-lint on project
.PHONY: lint
lint:
	@golangci-lint run .

## help: Show this help message
.PHONY: help
help: Makefile
	@echo "${APP_NAME}"
	@echo
	@echo " Choose a command run in "$(APP_NAME)":"
	@echo
	@sed -n 's/^## //p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo

# FUNCTIONS
define blue
	@tput setaf 4
	@echo $1
	@tput sgr0
endef