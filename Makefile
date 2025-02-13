.PHONY: api-lint
api-lint:
	buf lint

.PHONY: api-generate
api-generate:
	buf generate

.PHONY: test
test:
	go test -v -race -cover -timeout 10s ./...
