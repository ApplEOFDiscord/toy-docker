package main

import (
	"os"

	log "github.com/sirupsen/logrus"
	"gopkg.in/urfave/cli.v1"
)

const usage = "Just for fun!"

func main() {
	app := cli.NewApp()
	app.Name = "toy-docker"
	app.Usage = usage

	app.Commands = []cli.Command{
		initCommand,
		runCommand,
	}

	app.Before = func(context *cli.Context) error {
		log.SetFormatter(&log.JSONFormatter{})

		log.SetOutput(os.Stdout)
		return nil
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
