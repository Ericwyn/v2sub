package command

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

type CmdRunCb func(string)

func RunSyncForResultCb(cmdResCb CmdRunCb, name string, arg ...string) error {
	cmd := exec.Command(name, arg...)
	stdout, err := cmd.StdoutPipe()
	cmd.Stderr = cmd.Stdout

	if err != nil {
		return err
	}

	if err = cmd.Start(); err != nil {
		return err
	}
	for {
		tmp := make([]byte, 1024)
		_, err := stdout.Read(tmp)
		cmdResCb(string(tmp))
		if err != nil {
			break
		}
	}

	if err = cmd.Wait(); err != nil {
		return err
	}
	return nil
}

func RunSync(name string, arg ...string) error {
	cmd := exec.Command(name, arg...)
	stdout, err := cmd.StdoutPipe()
	cmd.Stderr = cmd.Stdout

	if err != nil {
		return err
	}

	if err = cmd.Start(); err != nil {
		return err
	}
	for {
		tmp := make([]byte, 1024)
		_, err := stdout.Read(tmp)
		fmt.Print(string(tmp))
		if err != nil {
			break
		}
	}

	if err = cmd.Wait(); err != nil {
		return err
	}
	return nil
}

func RunResult(s string) (string, error) {
	cmd := exec.Command("/bin/bash", "-c", s)

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()

	return out.String(), err
}

func StartBackground(name string, arg ...string) (*exec.Cmd, error) {
	cmd := exec.Command(name, arg...)

	devNull, err := os.OpenFile(os.DevNull, os.O_RDWR, 0)
	if err != nil {
		return nil, err
	}

	cmd.Stdout = devNull
	cmd.Stderr = devNull
	cmd.Stdin = devNull

	if err = cmd.Start(); err != nil {
		_ = devNull.Close()
		return nil, err
	}

	return cmd, nil
}

func FindProcess(pid int) (*os.Process, error) {
	return os.FindProcess(pid)
}

func IsProcessAlive(pid int) bool {
	if pid <= 0 {
		return false
	}
	process, err := FindProcess(pid)
	if err != nil {
		return false
	}
	if err = process.Signal(syscall.Signal(0)); err != nil {
		return false
	}
	return true
}

func TerminateProcess(pid int) error {
	process, err := FindProcess(pid)
	if err != nil {
		return err
	}
	return process.Signal(syscall.SIGTERM)
}

func KillProcess(pid int) error {
	process, err := FindProcess(pid)
	if err != nil {
		return err
	}
	return process.Signal(syscall.SIGKILL)
}
