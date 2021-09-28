.PHONY: all
all: fmt test swag doctoc

.PHONY: lint
lint:
	@echo "[golangci-lint] Running golangci-lint..."
	@golangci-lint run 2>&1
	@echo "------------------------------------[Done]"

.PHONY: fmt
fmt:
	@echo "[fmt] Formatting go project..."
	@gofmt -s -w . 2>&1
	@echo "------------------------------------[Done]"

.PHONY: test
test:
	@echo "[test] Running go test..."
	@go test ./... -coverprofile coverage.txt 2>&1
	@go tool cover -html=coverage.txt
	@echo "------------------------------------[Done]"

.PHONY: swag
swag:
	@echo "[swag] Running swag... (Looking up at internal/ and output to ./docs)"
	@swag init --parseInternal --propertyStrategy camelcase
	@echo "------------------------------------[Done]"

.PHONY: doctoc
doctoc:
	@echo "[doctoc] Running doctoc..."
	@doctoc . 2>&1
	@echo "------------------------------------[Done]"
