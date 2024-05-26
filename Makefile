PHONY: all
all:
	echo "Hello World"


PHONY: generate
generate:
	@mkdir -p generated/go
	@docker run --rm -u 1000 -v ${HOST_PWD}:/local \
		openapitools/openapi-generator-cli generate \
		-i /local/api/openapi.yaml \
		-g go \
		-o /local/generated/go
	@echo "✅ Generation complete"
