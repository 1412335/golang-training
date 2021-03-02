package main

import (
	"fw/audit/handler"
	pb "fw/audit/proto"
	"fw/pkg/broker"
	pkgConfig "fw/pkg/config"

	"github.com/micro/micro/v3/service"
	microBroker "github.com/micro/micro/v3/service/broker"
	"github.com/micro/micro/v3/service/logger"

	"context"
	"fw/audit/config"
	"fw/pkg/audit"
	"fw/pkg/dal/postgres"
	"os"
)

func loadConfig() *config.ServiceConfig {
	cfgFile := os.Getenv("CONFIG_FILE")
	srvConfigs := &config.ServiceConfig{}
	if err := pkgConfig.LoadConfig(cfgFile, srvConfigs); err != nil {
		logger.Fatalf("Load config failed: %v", err)
	}
	logger.Infof("Load config success: %v \n %+v", cfgFile, srvConfigs.Database)
	return srvConfigs
}

func main() {
	// Create service
	srv := service.New(
		service.Name("audit"),
		service.Version("latest"),
	)

	srv.Init()

	// load config
	srvConfigs := loadConfig()

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
		&handler.Audit{},
	); err != nil {
		logger.Fatalf("Error migrate database: %v", err)
	}

	// setup nats broker
	broker := broker.New(microBroker.DefaultBroker)
	defer broker.Disconnect()

	srvHandler := handler.New(db, broker)
	if err := srvHandler.SubscribeMessage(audit.AuditTopic, audit.AuditQueueInsert); err != nil {
		logger.Fatalf("sub broker failed: %v", err)
	}

	// Register handler
	if err := pb.RegisterAuditHandler(srv.Server(), srvHandler); err != nil {
		logger.Fatalf("register failed: %v", err)
	}

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
