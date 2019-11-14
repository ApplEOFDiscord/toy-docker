package container

import (
	"os"
	"syscall"

	log "github.com/sirupsen/logrus"
)

func RunContainerInitProcess(command string, args []string) error {
	log.Infof("In RunContainerInitProcess")

	//mount namespaces are shared by default, so we should explicitly declare them to be independent
	syscall.Mount("", "/", "", syscall.MS_PRIVATE|syscall.MS_REC, "")

	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")
	argv := []string{command}

	if err := syscall.Exec(command, argv, os.Environ()); err != nil {
		log.Errorf(err.Error())
		return err
	}
	return nil
}
