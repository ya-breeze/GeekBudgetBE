.PHONY: all
all: run

.PHONY: build
build:
	@go build -o bin/geekbudget cmd/main.go

.PHONY: run
run:
	@go run cmd/main.go server

.PHONY: replace-templates
replace-templates:
	@rm -rf pkg/generated/templates/goclient pkg/generated/templates/goserver
	@mkdir -p pkg/generated/templates/goclient pkg/generated/templates/goserver
	@docker run --rm -u 1000 -v ${HOST_PWD}:/local \
		openapitools/openapi-generator-cli author template -g go \
		-o /local/pkg/generated/templates/goclient
	@docker run --rm -u 1000 -v ${HOST_PWD}:/local \
		openapitools/openapi-generator-cli author template -g go-server \
		-o /local/pkg/generated/templates/goserver

.PHONY: generate
generate:
	@rm -rf pkg/generated/goclient pkg/generated/goserver
	@mkdir -p pkg/generated/goclient pkg/generated/goserver
	@docker run --rm -u 1000 -v ${HOST_PWD}:/local \
		openapitools/openapi-generator-cli generate \
		-i /local/api/openapi.yaml \
		-g go \
		-t /local/pkg/generated/templates/goclient \
		-o /local/pkg/generated/goclient \
		--additional-properties=packageName=goclient,withGoMod=false
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
		-t /local/pkg/generated/templates/goserver \
		-o /local/pkg/generated/goserver \
		--additional-properties=packageName=goserver,featureCORS=true,hideGenerationTimestamp=true
	@rm -rf \
		pkg/generated/goserver/api \
		pkg/generated/goserver/.openapi-generator-ignore \
		pkg/generated/goserver/Dockerfile \
		pkg/generated/goserver/go.*
	@mv -f pkg/generated/goserver/go/* pkg/generated/goserver
	@rm -rf pkg/generated/goserver/go
	@goimports -l -w ./pkg/generated/
	@gofumpt -l -w ./pkg/generated/
	@echo "✅ Generation complete"

.PHONY: validate
validate:
	@docker run --rm -v ${HOST_PWD}:/local openapitools/openapi-generator-cli validate -i /local/api/openapi.yaml
	@echo "✅ Validation complete"

.PHONY: lint
lint: validate
	@golangci-lint run
	@gofumpt -l -d .
	@echo "✅ Lint complete"

.PHONY: test
test:
	@ginkgo -r
	@echo "✅ Tests complete"