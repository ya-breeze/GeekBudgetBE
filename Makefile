PHONY: all
all:
	echo "Hello World"


PHONY: generate
generate:
	@rm -rf generated
	@mkdir -p generated
	@docker run --rm -u 1000 -v ${HOST_PWD}:/local \
		openapitools/openapi-generator-cli generate \
		-i /local/api/openapi.yaml \
		-g go \
		-o /local/generated/goclient
	@docker run --rm -u 1000 -v ${HOST_PWD}:/local \
		openapitools/openapi-generator-cli generate \
		-i /local/api/openapi.yaml \
		-g go-server \
		-o /local/generated/goserver
	@echo "✅ Generation complete"
