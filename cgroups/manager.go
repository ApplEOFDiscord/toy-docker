package cgroups

import (
	"toy-docker/cgroups/subsystems"

	log "github.com/sirupsen/logrus"
)

type CgroupManager struct {
	Path     string
	Resource *subsystems.ResourceConfig
}

func NewCgroupManager(path string) *CgroupManager {
	return &CgroupManager{
		Path: path,
	}
}

func (c *CgroupManager) Apply(pid int) error {
	for _, subsysIns := range subsystems.SubsystemIns {
		subsysIns.Apply(c.Path, pid)
	}
	return nil
}

func (c *CgroupManager) Set(res *subsystems.ResourceConfig) error {
	for _, subsysIns := range subsystems.SubsystemIns {
		subsysIns.Set(c.Path, res)
	}
	return nil
}

func (c *CgroupManager) Destroy() error {
	for _, subsysIns := range subsystems.SubsystemIns {
		if err := subsysIns.Remove(c.Path); err != nil {
			log.Warnf("Error remove cgroup %v", err)
		}
	}
	return nil
}
