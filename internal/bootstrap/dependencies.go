package bootstrap

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"statehound/internal/command"
	"statehound/internal/logger"
	"statehound/internal/model"
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

	var err error
	switch distro {
	case distroFedora:
		err = command.Run("dnf", append([]string{"install", "-y"}, packages...)...)
	case distroDebian:
		_ = command.Run("apt-get", "update", "-qq")
		err = command.Run("apt-get", append([]string{"install", "-y"}, packages...)...)
	case distroArch:
		err = command.Run("pacman", append([]string{"-S", "--noconfirm"}, packages...)...)
	case distroSUSE:
		err = command.Run("zypper", append([]string{"install", "-y"}, packages...)...)
	}

	if err != nil {
		return err
	}

	saveInstalledDeps(distro, packages)
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

type installedDeps struct {
	Packages []string `json:"packages"`
	Distro   string   `json:"distro"`
}

func saveInstalledDeps(distro distroFamily, packages []string) {
	data := installedDeps{
		Packages: packages,
		Distro:   distroName(distro),
	}

	b, err := json.Marshal(data)
	if err != nil {
		return
	}

	_ = os.WriteFile(model.InstalledDepsPath, b, 0600)
}

func distroName(distro distroFamily) string {
	switch distro {
	case distroFedora:
		return "fedora"
	case distroDebian:
		return "debian"
	case distroArch:
		return "arch"
	case distroSUSE:
		return "suse"
	default:
		return "unknown"
	}
}

func RemoveInstalledDeps() {
	data, err := os.ReadFile(model.InstalledDepsPath)
	if err != nil {
		return
	}

	var deps installedDeps
	if err := json.Unmarshal(data, &deps); err != nil {
		return
	}

	if len(deps.Packages) == 0 {
		return
	}

	logger.Status("removing dependencies installed by statehound: " + strings.Join(deps.Packages, ", "))

	switch deps.Distro {
	case "fedora":
		_ = command.Run("dnf", append([]string{"remove", "-y"}, deps.Packages...)...)
	case "debian":
		_ = command.Run("apt-get", append([]string{"remove", "-y"}, deps.Packages...)...)
	case "arch":
		_ = command.Run("pacman", append([]string{"-R", "--noconfirm"}, deps.Packages...)...)
	case "suse":
		_ = command.Run("zypper", append([]string{"remove", "-y"}, deps.Packages...)...)
	}

	_ = os.Remove(model.InstalledDepsPath)
}
