# all: default

.PHONY: all
all:
	docker-compose up

.PHONY: formatter
formatter:
	# go run ./pkg/cmd/formatter.go
	go run ./formatter/formatter.go

.PHONY: publisher
publisher:
	go run ./publisher/publisher.go

.PHONY: default
default:
	go run main.go nn


.PHONY: backend
backend:
	docker run \
		-d \
	 	--rm \
		--name jaeger \
		-p6831:6831/udp \
		-p16686:16686 \
		jaegertracing/all-in-one:latest

.PHONY: up
up:
	docker-compose -f ./docker/jaeger-elasticsearch_2.yml up -d

.PHONY: ps
ps:
	docker-compose -f ./docker/jaeger-elasticsearch_2.yml ps

.PHONY: down
down:
	docker-compose -f ./docker/jaeger-elasticsearch_2.yml down

.PHONY: genproto
genproto:
	protoc.exe -I ./services/publisher/proto/ --go_out=plugins=grpc:./services/publisher ./services/publisher/proto/*.proto

.PHONY: test
test:
	curl -v 'localhost:8080/format?helloTo=luyen&greeting=love&num=3'

.PHONY: grafana
grafana:
	docker-compose -f ./docker/grafana-integration/docker-compose.yml up


# https://github.com/spf13/cast/blob/master/Makefile
GOVERSION := $(shell go version | cut -d ' ' -f 3 | cut -d '.' -f 2)

# .PHONY: check fmt lint test test-race vet test-cover-html help
# .DEFAULT_GOAL := help

# check: test-race fmt vet lint ## Run tests and linters

# test: ## Run tests
# 	go test ./...

# test-race: ## Run tests with race detector
# 	go test -race ./...

fmt: ## Run gofmt linter
ifeq "$(GOVERSION)" "12"
	@for d in `go list` ; do \
		if [ "`gofmt -l -s $$GOPATH/src/$$d | tee /dev/stderr`" ]; then \
			echo "^ improperly formatted go files" && echo && exit 1; \
		fi \
	done
endif

lint: ## Run golint linter
	@for d in `go list` ; do \
		if [ "`golint $$d | tee /dev/stderr`" ]; then \
			echo "^ golint errors!" && echo && exit 1; \
		fi \
	done

vet: ## Run go vet linter
	@if [ "`go vet | tee /dev/stderr`" ]; then \
		echo "^ go vet errors!" && echo && exit 1; \
	fi

# test-cover-html: ## Generate test coverage report
# 	go test -coverprofile=coverage.out -covermode=count
# 	go tool cover -func=coverage.out

# help:
# 	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'