package daemon

import (
	"fmt"
	"os"
	"os/user"
	"strconv"

	"statehound/internal/model"
)

func prepareRuntimeDir() error {
	group, err := user.LookupGroup(model.GroupName)
	if err != nil {
		return fmt.Errorf("failed to lookup group %s: %w", model.GroupName, err)
	}

	gid, err := strconv.Atoi(group.Gid)
	if err != nil {
		return fmt.Errorf("invalid group gid %q: %w", group.Gid, err)
	}

	if err := os.Chown(model.RuntimeDir, 0, gid); err != nil {
		return err
	}

	if err := os.Chmod(model.RuntimeDir, 0750); err != nil {
		return err
	}

	return nil
}
