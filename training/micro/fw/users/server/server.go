package main

import (
	"context"
	"fw/pkg/audit"
	"fw/pkg/broker"
	pkgConfig "fw/pkg/config"
	"fw/pkg/dal/postgres"
	"fw/users/config"
	"fw/users/handler"
	pb "fw/users/proto"
	"os"
	// "time"
	// "path/filepath"
	"strings"

	"github.com/micro/micro/v3/service"
	// microConfig "github.com/micro/micro/v3/service/config"
	microBroker "github.com/micro/micro/v3/service/broker"
	"github.com/micro/micro/v3/service/context/metadata"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/server"
	// "github.com/asim/go-micro/v3/config"
	// // "github.com/asim/go-micro/v3/config/source"
	// "github.com/asim/go-micro/v3/config/source/env"
	// "github.com/asim/go-micro/v3/config/source/file"
	// "gorm.io/driver/postgres"
	// "gorm.io/gorm"
)

type Authentication struct {
	jwtManager *handler.JWTManager
}

func (a *Authentication) AuthWrapper(fn server.HandlerFunc) server.HandlerFunc {
	return func(ctx context.Context, req server.Request, rsp interface{}) error {
		// User login or create
		if req.Endpoint() == "Users.Create" || req.Endpoint() == "Users.Auth" {
			// TEST
			ctx2 := metadata.Set(ctx, "UserID", "1111")
			return fn(ctx2, req, rsp)
		}

		// read authorization token from context metadata
		meta, ok := metadata.FromContext(ctx)
		if !ok {
			return errors.BadRequest("METADATA_NIL", "")
		}

		auth, ok := meta["Authorization"]
		if !ok {
			return errors.Unauthorized("AUTH_NIL", "")
		}
		authSplit := strings.Split(auth, " ")
		// token should be "Bearer ..."
		if len(authSplit) != 2 {
			return errors.Unauthorized("AUTH_INCORRECT", "")
		}

		token := authSplit[1]

		logger.Infof("Endpoint: %v", req.Endpoint())

		// validate token
		claims, err := a.jwtManager.Verify(token)
		if err != nil {
			return errors.Unauthorized("AUTH_INCORRECT", "")
		}

		// Add current user to context to use in saving audit records
		ctx2 := metadata.Set(ctx, "UserID", claims.ID)

		return fn(ctx2, req, rsp)
	}
}

func loadConfig() *config.ServiceConfig {
	// conf, err := config.NewConfig()
	// if err != nil {
	// 	logger.Fatalf("Expected no error but got %v", err)
	// }
	// path := "./config.json"
	// if err := conf.Load(
	// 	file.NewSource(
	// 		file.WithPath(path),
	// 	),
	// 	env.NewSource(),
	// ); err != nil {
	// 	logger.Fatalf("Expected no error but got %v", err)
	// }

	// actualHost := conf.Get("amqp", "host").String("backup")
	// if actualHost != "rabbit.testing.com" {
	// 	logger.Fatalf("Expected %v but got %v",
	// 		"rabbit.testing.com",
	// 		actualHost)
	// }

	cfgFile := os.Getenv("CONFIG_FILE")
	srvConfigs := &config.ServiceConfig{}
	if err := pkgConfig.LoadConfig(cfgFile, srvConfigs); err != nil {
		logger.Fatalf("Load config failed: %v", err)
	}
	logger.Infof("Load config success: %v \n %+v \n %+v", cfgFile, srvConfigs.JWT, srvConfigs.Database)
	return srvConfigs
}

func main() {
	// load config
	srvConfigs := loadConfig()
	// jwtManager
	jwtManager := handler.NewJWTManager(srvConfigs.JWT)
	//
	auth := Authentication{
		jwtManager: jwtManager,
	}

	// Create service
	srv := service.New(
		service.Name("users"),
		service.Version("latest"),
		service.WrapHandler(auth.AuthWrapper),
	)

	// optionally setup command line usage
	srv.Init()

	dal, err := postgres.NewDataAccessLayer(context.Background(), srvConfigs.Database)
	if err != nil {
		logger.Fatalf("Error connect database: %v", err)
	}
	// defer dal.Disconnect()
	db := dal.GetDatabase()
	if db == nil {
		logger.Fatalf("Error connect database: %v", err)
	}
	if err := db.AutoMigrate(
		&handler.User{},
	); err != nil {
		logger.Fatalf("Error migrate database: %v", err)
	}

	// setup nats broker
	broker := broker.New(microBroker.DefaultBroker)
	defer broker.Disconnect()

	var handlerOpts []handler.UsersHandlerOption
	// audit
	if srvConfigs.EnableAuditRecords {
		handlerOpts = append(handlerOpts, handler.WithAudit(&audit.Audit{Broker: broker}))
	}

	usersHandler, err := handler.NewUsersHandler(db, jwtManager, handlerOpts...)
	if err != nil {
		logger.Fatalf("Error new handler: %v", err)
	}

	// Register handler
	if err := pb.RegisterUsersHandler(srv.Server(), usersHandler); err != nil {
		logger.Fatalf("Error registering handler: %v", err)
	}

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
