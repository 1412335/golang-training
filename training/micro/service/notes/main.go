package main

import (
	"notes/handler"
	pb "notes/proto"

	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("notes"),
		service.Version("latest"),
	)

	// Register handler
	pb.RegisterNotesHandler(srv.Server(), new(handler.Notes))

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
