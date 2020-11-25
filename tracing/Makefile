all: default

.PHONY: formatter
formatter:
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
