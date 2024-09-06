ROOT_DIR := $(dir $(realpath $(lastword $(MAKEFILE_LIST))))

.PHONY: all
all: build test lint
	@echo "🎉 You are good to go!"


.PHONY: build
build:
	@cd cmd && go build -o ../bin/geekbudget
	@echo "✅ Build complete"

.PHONY: run
run: build
	@GB_USERS=test:JDJhJDEwJC9sVWJpTlBYVlZvcU9ZNUxIZmhqYi4vUnRuVkJNaEw4MTQ2VUdFSXRDeE9Ib0ZoVkRLR3pl \
	GB_PREFILL=true \
	GB_DBPATH=$(ROOT_DIR)geekbudget.db \
	./bin/geekbudget server

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
	# Golang client and server
	@rm -rf pkg/generated/goclient pkg/generated/goserver pkg/generated/angular
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

	# Angular client
	@docker run --rm -u 1000 -v ${HOST_PWD}:/local \
		openapitools/openapi-generator-cli generate \
		-i /local/api/openapi.yaml \
		-g typescript-angular \
		-o /local/pkg/generated/angular \
		--additional-properties=apiModulePrefix=GeekbudgetClient,configurationPrefix=GeekbudgetClient

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

.PHONY: watch
watch:
	@ginkgo watch -r
