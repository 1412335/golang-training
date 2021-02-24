
GOPATH:=$(shell go env GOPATH)
.PHONY: init
init:
	go get -u github.com/golang/protobuf/proto
	go get -u github.com/golang/protobuf/protoc-gen-go
	go get github.com/micro/micro/v3/cmd/protoc-gen-micro
.PHONY: proto
proto:
	protoc --proto_path=. --micro_out=. --go_out=:. proto/tags.proto
	
.PHONY: build
build:
	go build -o tags *.go

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: docker
docker:
	docker build . -t tags:latest

.PHONY: run
run: 
	micro run .

.PHONY: update
update: 
	micro update .

.PHONY: logs
logs: 
	micro logs -f tags