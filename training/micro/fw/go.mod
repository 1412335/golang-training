module fw

go 1.15

require (
	github.com/1412335/grpc-rest-microservice v0.0.0-20201124041515-eb8e377d77f1 // indirect
	github.com/asim/go-micro/v3 v3.5.0 // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/golang/protobuf v1.4.3
	github.com/google/uuid v1.1.2
	github.com/micro/micro/v3 v3.0.0
	github.com/spf13/viper v1.7.1
	golang.org/x/crypto v0.0.0-20200820211705-5c72a883971a
	google.golang.org/protobuf v1.25.0
	gorm.io/driver/postgres v1.0.8
	gorm.io/gorm v1.20.12
)

// This can be removed once etcd becomes go gettable, version 3.4 and 3.5 is not,
// see https://github.com/etcd-io/etcd/issues/11154 and https://github.com/etcd-io/etcd/issues/11931.
replace google.golang.org/grpc => google.golang.org/grpc v1.26.0
