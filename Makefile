.PHONY: api-lint
api-lint:
	buf lint

.PHONY: api-generate
api-generate:
	buf generate
