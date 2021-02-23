package main

import (
	"context"
	"fmt"
	"helloworld/handler"
	pb "helloworld/proto"
	"os"
	"time"

	// "github.com/micro/micro/v3/cmd"
	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"
)

func runClient() {
	// create and initialise a new service
	srv := service.New()

	// create the proto client for helloworld
	client := pb.NewHelloworldService("helloworld", srv.Client())

	// call an endpoint on the service
	rsp, err := client.Call(context.Background(), &pb.Request{
		Name: "John",
	})
	if err != nil {
		fmt.Println("Error calling helloworld: ", err)
		return
	}

	// print the response
	fmt.Println("Response: ", rsp.Msg)

	// let's delay the process for exiting for reasons you'll see below
	time.Sleep(time.Second * 5)
}

func main() {
	// Create service
	srv := service.New(
		service.Name("helloworld"),
		service.Version("latest"),
		service.Metadata(map[string]string{
			"type": "helloworld",
		}),
		// Setup some flags. Specify --run_client to run the client
		// Add runtime flags
		// We could do this below too
		// service.Flags(&cli.BoolFlag{
		// 	Name:  "run_client",
		// 	Usage: "run client",
		// }),
	)

	// Init will parse the command line flags. Any flags set will
	// override the above settings. Options defined here will
	// override anything set on the command line.
	// srv.Init(
	// 	service.Action(func(c *cli.Context) {
	// 		if c.Bool("run_client") {
	// 			runClient()
	// 			os.Exit(0)
	// 		}
	// 		return nil
	// 	}),
	// )

	// Register handler
	pb.RegisterHelloworldHandler(srv.Server(), new(handler.Helloworld))

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
