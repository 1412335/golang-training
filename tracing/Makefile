all: default

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

.PHONY: genproto
genproto:
	protoc.exe -I ./services/publisher/proto/ --go_out=plugins=grpc:./services/publisher ./services/publisher/proto/*.proto