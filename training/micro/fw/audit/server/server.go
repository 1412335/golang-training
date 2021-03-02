package main

import (
	"fw/audit/handler"
	pb "fw/audit/proto"
	"fw/pkg/broker"

	"github.com/micro/micro/v3/service"
	microBroker "github.com/micro/micro/v3/service/broker"
	"github.com/micro/micro/v3/service/logger"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("audit"),
		service.Version("latest"),
	)

	srv.Init()

	// setup nats broker
	broker := broker.New(microBroker.DefaultBroker)
	defer broker.Disconnect()

	err := broker.SubMsg("Audit", "Audit.Queue")
	if err != nil {
		logger.Fatalf("sub broker failed: %v", err)
	}

	// Register handler
	pb.RegisterAuditHandler(srv.Server(), new(handler.Audit))

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
