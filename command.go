package main

import (
	"fmt"
	"toy-docker/container"

	log "github.com/sirupsen/logrus"
	"gopkg.in/urfave/cli.v1"
)

var runCommand = cli.Command{
	Name:  "run",
	Usage: "Create a container with namespace and cgroups limit: toy-docker run -it [command]",
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "it",
			Usage: "enable tty",
		},
	},

	Action: func(context *cli.Context) error {
		if len(context.Args()) == 0 {
			return fmt.Errorf("Missing container command")
		}
		cmd := context.Args().Get(0)
		tty := context.Bool("it")
		Run(tty, cmd)
		return nil
	},
}

var initCommand = cli.Command{
	Name:  "init",
	Usage: "Init container process. Do not call it outside",

	Action: func(context *cli.Context) error {
		cmd := context.Args().Get(0)
		log.Infof("command %s", cmd)
		err := container.RunContainerInitProcess(cmd, nil)
		return err
	},
}
