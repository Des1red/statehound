package bootstrap

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"statehound/internal/command"
	"statehound/internal/logger"
)

type distroFamily int

const (
	distroUnknown distroFamily = iota
	distroFedora
	distroDebian
	distroArch
	distroSUSE
)

func detectDistro() distroFamily {
	f, err := os.Open("/etc/os-release")
	if err != nil {
		return distroUnknown
	}
	defer f.Close()

	values := make(map[string]string)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		values[key] = strings.Trim(value, `"`)
	}

	combined := strings.ToLower(values["ID"] + " " + values["ID_LIKE"])

	switch {
	case containsAny(combined, "fedora", "rhel", "centos", "almalinux", "rocky"):
		return distroFedora
	case containsAny(combined, "debian", "ubuntu", "kali", "mint", "pop"):
		return distroDebian
	case containsAny(combined, "arch", "manjaro", "endeavour"):
		return distroArch
	case containsAny(combined, "suse", "opensuse"):
		return distroSUSE
	default:
		return distroUnknown
	}
}

func containsAny(s string, needles ...string) bool {
	for _, needle := range needles {
		if strings.Contains(s, needle) {
			return true
		}
	}
	return false
}

func checkNotifierDeps() []string {
	var missing []string
	for _, dep := range []string{"notify-send", "setfacl"} {
		if _, err := command.Output("which", dep); err != nil {
			missing = append(missing, dep)
		}
	}
	return missing
}

var packageMap = map[distroFamily]map[string]string{
	distroFedora: {
		"notify-send": "libnotify",
		"setfacl":     "acl",
	},
	distroDebian: {
		"notify-send": "libnotify-bin",
		"setfacl":     "acl",
	},
	distroArch: {
		"notify-send": "libnotify",
		"setfacl":     "acl",
	},
	distroSUSE: {
		"notify-send": "libnotify-tools",
		"setfacl":     "acl",
	},
}

func installNotifierDeps(distro distroFamily, missing []string) error {
	distroPackages, ok := packageMap[distro]
	if !ok {
		return fmt.Errorf("unsupported distro — install manually: %s", strings.Join(missing, ", "))
	}

	var packages []string
	for _, binary := range missing {
		if pkg, ok := distroPackages[binary]; ok {
			packages = append(packages, pkg)
		}
	}

	if len(packages) == 0 {
		return nil
	}

	switch distro {
	case distroFedora:
		return command.Run("dnf", append([]string{"install", "-y"}, packages...)...)
	case distroDebian:
		_ = command.Run("apt-get", "update", "-qq")
		return command.Run("apt-get", append([]string{"install", "-y"}, packages...)...)
	case distroArch:
		return command.Run("pacman", append([]string{"-S", "--noconfirm"}, packages...)...)
	case distroSUSE:
		return command.Run("zypper", append([]string{"install", "-y"}, packages...)...)
	}

	return nil
}

func ensureNotifierDeps() {
	missing := checkNotifierDeps()
	if len(missing) == 0 {
		return
	}

	logger.Status("installing missing notifier dependencies: " + strings.Join(missing, ", "))

	distro := detectDistro()
	if distro == distroUnknown {
		logger.Warn("could not detect distro — install manually: " + strings.Join(missing, ", "))
		return
	}

	if err := installNotifierDeps(distro, missing); err != nil {
		logger.Warn("failed to install notifier dependencies: " + err.Error())
	}
}
