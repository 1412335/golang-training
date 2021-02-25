module chats

go 1.15

require (
	github.com/golang/protobuf v1.4.3
	github.com/google/uuid v1.1.2
	github.com/micro/micro/v3 v3.0.0
	github.com/stretchr/testify v1.7.0
	google.golang.org/protobuf v1.25.0
	gorm.io/driver/postgres v1.0.8
	gorm.io/gorm v1.20.12
)

// This can be removed once etcd becomes go gettable, version 3.4 and 3.5 is not,
// see https://github.com/etcd-io/etcd/issues/11154 and https://github.com/etcd-io/etcd/issues/11931.
replace google.golang.org/grpc => google.golang.org/grpc v1.26.0
