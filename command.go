package main

import (
	"fmt"
	"toy-docker/cgroups/subsystems"
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
		cli.StringFlag{
			Name:  "m",
			Usage: "memory limit",
		},
		cli.StringFlag{
			Name:  "cpuset",
			Usage: "cpuset limit",
		},
		cli.StringFlag{
			Name:  "cpushare",
			Usage: "cpushare limit",
		},
	},

	Action: func(context *cli.Context) error {
		if len(context.Args()) == 0 {
			return fmt.Errorf("Missing container command")
		}

		tty := context.Bool("it")
		var cmdArray []string
		for _, cmd := range context.Args() {
			cmdArray = append(cmdArray, cmd)
		}
		res := &subsystems.ResourceConfig{
			MemoryLimit: context.String("m"),
			CpuSet:      context.String("cpuset"),
			CpuShare:    context.String("cpushare"),
		}

		Run(tty, cmdArray, res)
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
