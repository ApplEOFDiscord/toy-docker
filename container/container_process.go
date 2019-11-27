package container

import (
	"os"
	"os/exec"
	"syscall"

	log "github.com/sirupsen/logrus"
)

func NewParentProcess(tty bool) (*exec.Cmd, *os.File) {
	readPipe, writePipe, err := NewPipe()
	if err != nil {
		log.Errorf("New pipe error %v", err)
		return nil, nil
	}

	cmd := exec.Command("/proc/self/exe", "init")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID |
			syscall.CLONE_NEWNS | syscall.CLONE_NEWNET | syscall.CLONE_NEWIPC,
	}
	if tty {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	cmd.ExtraFiles = []*os.File{readPipe}

	mntPath := "/root/mnt/"
	rootPath := "/root/"
	NewWorkSpace(rootPath, mntPath)

	cmd.Dir = mntPath
	return cmd, writePipe
}

func NewPipe() (*os.File, *os.File, error) {
	read, write, err := os.Pipe()
	if err != nil {
		return nil, nil, err
	}
	return read, write, nil
}

func NewWorkSpace(rootPath string, mntPath string) {
	CreateReadOnlyLayer(rootPath)
	CreateWriteLayer(rootPath)
	CreateMountPoint(rootPath, mntPath)
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func CreateReadOnlyLayer(rootPath string) {
	busyboxPath := rootPath + "busybox/"
	busyboxTarPath := rootPath + "busybox.tar"
	exist, err := PathExists(busyboxPath)
	if err != nil {
		log.Infof("Fail to judge whether dir %s exists %v", busyboxPath, err)
	}
	if exist == false {
		if err := os.Mkdir(busyboxPath, 0777); err != nil {
			log.Errorf("Mkdir dir %s error %v", busyboxPath, err)
		}
		if _, err := exec.Command("tar", "-xvf", busyboxTarPath, "-C", busyboxPath).CombinedOutput(); err != nil {
			log.Errorf("Untar dir %s error", busyboxTarPath, err)
		}
	}
}

func CreateWriteLayer(rootPath string) {
	writePath := rootPath + "write_layer/"
	if err := os.Mkdir(writePath, 0777); err != nil {
		log.Errorf("Mkdir dir %s error %v", writePath, err)
	}
}

func CreateMountPoint(rootPath string, mntPath string) {
	//Create mnt directory to be the mount point
	if err := os.Mkdir(mntPath, 0777); err != nil {
		log.Errorf("Mkdir dir %s error %v", mntPath, err)
	}

	dirs := "dirs=" + rootPath + "write_layer:" + rootPath + "busybox"
	cmd := exec.Command("mount", "-t", "aufs", "-o", dirs, "none", mntPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Errorf("Mount aufs error %v", err)
	}
}

func DeleteWorkSpace(rootPath string, mntPath string) {
	DeleteMountPoint(rootPath, mntPath)
	DeleteWriteLayer(rootPath)
}

func DeleteMountPoint(rootPath string, mntPath string) {
	cmd := exec.Command("umount", mntPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Errorf("%v", err)
	}
	if err := os.RemoveAll(mntPath); err != nil {
		log.Errorf("Remove dir %s error %v", mntPath, err)
	}
}

func DeleteWriteLayer(rootPath string) {
	writePath := rootPath + "write_layer/"
	if err := os.RemoveAll(writePath); err != nil {
		log.Errorf("Remove dir %s error %v", writePath, err)
	}
}
