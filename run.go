package main

import (
	"os"
	"strings"
	"toy-docker/cgroups"
	"toy-docker/cgroups/subsystems"
	"toy-docker/container"

	log "github.com/sirupsen/logrus"
)

func Run(tty bool, cmdArray []string, res *subsystems.ResourceConfig) {
	parent, writePipe := container.NewParentProcess(tty)
	if parent == nil {
		log.Errorf("New parent process error")
	}
	if err := parent.Start(); err != nil {
		log.Error(err)
	}

	cgroupManager := cgroups.NewCgroupManager("toy-docker-cgroup")
	defer cgroupManager.Destroy()

	cgroupManager.Set(res)
	cgroupManager.Apply(parent.Process.Pid)

	sendInitCommand(cmdArray, writePipe)
	parent.Wait()
	os.Exit(-1)
}

func sendInitCommand(cmdArray []string, writePipe *os.File) {
	command := strings.Join(cmdArray, " ")
	log.Infof("Command is %s", command)
	writePipe.WriteString(command)
	writePipe.Close()
}
