PHONY: all
all:
	echo "Hello World"


PHONY: generate
generate:
	@rm -rf pkg/generated
	@mkdir -p pkg/generated
	@docker run --rm -u 1000 -v ${HOST_PWD}:/local \
		openapitools/openapi-generator-cli generate \
		-i /local/api/openapi.yaml \
		-g go \
		-o /local/pkg/generated/goclient \
		--additional-properties=packageName=goclient
	@docker run --rm -u 1000 -v ${HOST_PWD}:/local \
		openapitools/openapi-generator-cli generate \
		-i /local/api/openapi.yaml \
		-g go-server \
		-o /local/pkg/generated/goserver
		--additional-properties=packageName=goserver
	@echo "✅ Generation complete"

PHONY: validate
validate:
	@docker run --rm -v ${HOST_PWD}:/local openapitools/openapi-generator-cli validate -i /local/api/openapi.yaml
	@echo "✅ Validation complete"

PHONY: lint
lint: validate
