.PHONY: all
all: run

.PHONY: build
build:
	@go build -o bin/geekbudget cmd/main.go

.PHONY: run
run:
	@go run cmd/main.go

.PHONY: generate
generate:
	@rm -rf pkg/generated
	@mkdir -p pkg/generated
	@docker run --rm -u 1000 -v ${HOST_PWD}:/local \
		openapitools/openapi-generator-cli generate \
		-i /local/api/openapi.yaml \
		-g go \
		-o /local/pkg/generated/goclient \
		--additional-properties=packageName=goclient
	@rm -rf \
		pkg/generated/goclient/api \
		pkg/generated/goclient/.gitignore \
		pkg/generated/goclient/.openapi-generator-ignore \
		pkg/generated/goclient/.travis.yml \
		pkg/generated/goclient/*.sh \
		pkg/generated/goclient/go.* \
		pkg/generated/goclient/test
	@docker run --rm -u 1000 -v ${HOST_PWD}:/local \
		openapitools/openapi-generator-cli generate \
		-i /local/api/openapi.yaml \
		-g go-server \
		-o /local/pkg/generated/goserver \
		--additional-properties=packageName=goserver
	@rm -rf \
		pkg/generated/goserver/api \
		pkg/generated/goserver/.openapi-generator-ignore \
		pkg/generated/goserver/Dockerfile \
		pkg/generated/goserver/go.*
	@mv -f pkg/generated/goserver/go/* pkg/generated/goserver
	@rm -rf pkg/generated/goserver/go
	@echo "✅ Generation complete"

.PHONY: validate
validate:
	@docker run --rm -v ${HOST_PWD}:/local openapitools/openapi-generator-cli validate -i /local/api/openapi.yaml
	@echo "✅ Validation complete"

.PHONY: lint
lint: validate
