package viewer

import (
	"fmt"
	"os"
	"os/exec"
)

func startNotebook(key int, activeTab int) (*exec.Cmd, error) {
	notebook := exec.Command(
		"yad",
		"--notebook",
		fmt.Sprintf("--key=%d", key),
		"--tab=Critical",
		"--tab=Normal",
		fmt.Sprintf("--active-tab=%d", activeTab),
		"--width=900",
		"--height=600",
		"--title=statehound events",
		"--no-buttons",
	)
	notebook.Env = buildEnv()
	notebook.Stderr = os.Stderr

	if err := notebook.Start(); err != nil {
		return nil, err
	}

	return notebook, nil
}
