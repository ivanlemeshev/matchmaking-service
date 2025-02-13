.PHONY: api-lint
api-lint:
	@buf lint

.PHONY: api-generate
api-generate:
	@buf generate

.PHONY: test
test:
	@go test -race -cover -timeout 30s ./...

.PHONY: run
run:
	go run cmd/main.go
