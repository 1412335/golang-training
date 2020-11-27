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
