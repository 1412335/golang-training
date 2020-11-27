package main

import (
	"golang-training/week2/command"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {

	// cli
	app := &cli.App{
		Commands: []*cli.Command{
			command.InfoCommand,
			command.MigrateCommand,
			command.ServeCommand,
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
