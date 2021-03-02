
GOPATH:=$(shell go env GOPATH)
.PHONY: init
init:
	go get -u github.com/golang/protobuf/proto
	go get -u github.com/golang/protobuf/protoc-gen-go
	go get github.com/micro/micro/v3/cmd/protoc-gen-micro

# gen proto
.PHONY: proto
proto:
	protoc --proto_path=$$GOPATH/src:. --micro_out=source_relative:.. --go_out=. --go_opt=paths=source_relative users/proto/users.proto
	
# format
.PHONY: fmt
fmt:
	go fmt ./...

# linter
.PHONY: lint
lint: fmt
	golangci-lint run ./...

# test
.PHONY: test
test: lint
	go test -v ./... -cover


# micro command
# start micro server
microserve:
	micro server

# Login to micro server
micrologin:
	micro login --username admin --password micro

# running single service
.PHONY: microstartusersrv
microstartusersrv:
	micro run --env_vars "CONFIG_FILE=$$PWD/config.yml,MICRO_BROKER=nats,MICRO_BROKER_ADDRESS=natsUser:TestDB@home2@localhost" --name users users/server

.PHONY: micrologsusersrv
micrologsusersrv:
	micro logs -f users

.PHONY: microstartsrvs
microstartsrvs:
	# make microkillsrvs
	make microstartusersrv

# Run the whole application locally, micro server must already be running locally
microlocalup:
	#micro server
	docker-compose up -d nats
	sleep 20s
	docker-compose ps
	make micrologin
	make microstartsrvs

microlocaldown:
	make microkillsrvs
	# docker-compose down

.PHONY: microkillsrvs
microkillsrvs:
	micro kill users
	
.PHONY: microupdatesrvs
microupdatesrvs:
	micro update users

# Docker
# microstartsrvscont:
# 	micro run --env_vars 'CONFIG_FILE=./config.yml' --name users users/server

## Start the whole application running in Docker
microup:
	docker-compose up -d microserver
	docker-compose up -d timescaledb nats
	sleep 20s
	docker exec microserver make micrologin
	docker exec microserver make microstartsrvs
	# docker-compose up web

# 
.PHONY: build
build:
	go build -o fw *.go

.PHONY: docker
docker:
	docker build . -t fw:latest


demo:
	echo "a=$$PWD"