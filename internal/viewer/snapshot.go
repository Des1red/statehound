package viewer

import (
	"fmt"
	"os/exec"

	"statehound/internal/logger"
)

func ShowSnapshot(fetch func() (string, error)) {
	for {
		content, err := fetch()
		if err != nil {
			logger.Failed("snapshot gui: failed to fetch snapshot", err)
			return
		}

		cmd := exec.Command(
			"yad",
			"--text-info",
			"--title=statehound snapshot",
			"--width=900",
			"--height=700",
			"--fontname=monospace",
			"--button=Refresh:1",
			"--button=Close:0",
		)
		cmd.Env = buildEnv()

		stdin, err := cmd.StdinPipe()
		if err != nil {
			logger.Failed("snapshot gui: failed to create pipe", err)
			return
		}

		if err := cmd.Start(); err != nil {
			logger.Failed("snapshot gui: failed to start yad", err)
			return
		}

		fmt.Fprintln(stdin, content)
		stdin.Close()
		cmd.Wait()

		if cmd.ProcessState.ExitCode() != 1 {
			return
		}
	}
}
