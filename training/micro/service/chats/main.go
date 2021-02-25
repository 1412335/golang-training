package main

import (
	"chats/handler"
	pb "chats/proto"

	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/config"
	"github.com/micro/micro/v3/service/logger"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const dbAddress = "postgresql://root:root@localhost:5432/chats?sslmode=disable"

func main() {
	// Create service
	srv := service.New(
		service.Name("chats"),
		service.Version("latest"),
	)

	// get config
	cfg, err := config.Get("chats.db")
	if err != nil {
		logger.Fatalf("Could not get config: %v", err)
	}
	// connect db
	dsn := cfg.String(dbAddress)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Fatalf("Could not connect to database: %v", err)
	}
	if err := db.AutoMigrate(&handler.Chat{}); err != nil {
		logger.Fatalf("Auto migrate failed: %v", err)
	}

	// Register handler
	if err := pb.RegisterChatsHandler(srv.Server(), &handler.Chats{
		DB: db,
	}); err != nil {
		logger.Fatalf("Could not register handler: %v", err)
	}

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
