package main

import (
	"posts/handler"
	pb "posts/proto"

	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"
	// tags "tags/proto"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("posts"),
		service.Version("latest"),
	)

	// Register handler
	pb.RegisterPostsHandler(srv.Server(), handler.NewPosts(
	// tags.NewTagsService("tags", srv.Client()),
	))
	// srv.Handler(handler.NewPosts())

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
