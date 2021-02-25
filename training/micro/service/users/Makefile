
GOPATH:=$(shell go env GOPATH)
.PHONY: init
init:
	go get -u github.com/golang/protobuf/proto
	go get -u github.com/golang/protobuf/protoc-gen-go
	go get github.com/micro/micro/v3/cmd/protoc-gen-micro
.PHONY: proto
proto:
	protoc --proto_path=. --micro_out=. --go_out=:. proto/users.proto
	
.PHONY: build
build:
	go build -o users *.go

# gofmt
.PHONY: fmt
fmt:
	go fmt ./...

# go-lint
.PHONY: lint
lint: fmt
	golangci-lint run ./...

.PHONY: test
test: lint
	go test -v ./... -cover

.PHONY: docker
docker:
	docker build . -t users:latest

.PHONY: status
status: 
	micro status

.PHONY: run
run: 
	micro run .

.PHONY: update
update: 
	micro update .

.PHONY: logs
logs: 
	micro logs -f users


.PHONY: db
db:
	docker exec -ti postgresql psql