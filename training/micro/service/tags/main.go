package main

import (
	"tags/handler"
	pb "tags/proto"

	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("tags"),
		service.Version("latest"),
	)

	// Register handler
	pb.RegisterTagsHandler(srv.Server(), handler.NewTags())

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
