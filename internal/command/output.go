package command

import "os/exec"

func Output(name string, args ...string) ([]byte, error) {
	return exec.Command(name, args...).Output()
}
