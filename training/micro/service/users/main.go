package main

import (
	"users/handler"
	pb "users/proto"

	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/config"
	"github.com/micro/micro/v3/service/logger"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var dbAddress = "postgresql://root:root@localhost:5432/users?sslmode=disable"

func main() {
	// Create service
	srv := service.New(
		service.Name("users"),
		service.Version("latest"),
	)

	// connect db
	cfg, err := config.Get("users.db")
	if err != nil {
		logger.Fatalf("Error loading config: %v", err)
	}
	dsn := cfg.String(dbAddress)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Fatalf("Error connect database: %v", err)
	}
	if err := db.AutoMigrate(&handler.User{}); err != nil {
		logger.Fatalf("Error migrate database: %v", err)
	}

	// Register handler
	pb.RegisterUsersHandler(srv.Server(), &handler.Users{
		DB: db,
	})

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
