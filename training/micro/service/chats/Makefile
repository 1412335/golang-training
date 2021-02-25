
GOPATH:=$(shell go env GOPATH)
.PHONY: init
init:
	go get -u github.com/golang/protobuf/proto
	go get -u github.com/golang/protobuf/protoc-gen-go
	go get github.com/micro/micro/v3/cmd/protoc-gen-micro
.PHONY: proto
proto:
	protoc --proto_path=. --micro_out=. --go_out=:. proto/chats.proto
	
.PHONY: build
build:
	go build -o chats *.go

# gofmt
.PHONY: fmt
fmt:
	# go fmt  -mod=readonly ./...
	go fmt ./...

# go-lint
.PHONY: lint
lint: fmt
	golangci-lint run ./...

.PHONY: test
test: lint
	go test -v ./... -cover

.PHONY: utest
utest: lint update
	go test -v ./... -cover

.PHONY: docker
docker:
	docker build . -t chats:latest

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
	micro logs -f chats


.PHONY: db
db:
	docker exec -ti postgresql psql