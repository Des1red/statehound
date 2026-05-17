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
	"statehound/internal/system"
)

type distroFamily int

const (
	distroUnknown distroFamily = iota
	distroFedora
	distroDebian
	distroArch
	distroSUSE
)

type depInfo struct {
	CheckCmd []string
	Packages map[distroFamily]string
}

var deps = map[string]depInfo{
	"notify-send": {
		CheckCmd: []string{"which", "notify-send"},
		Packages: map[distroFamily]string{
			distroFedora: "libnotify",
			distroDebian: "libnotify-bin",
			distroArch:   "libnotify",
			distroSUSE:   "libnotify-tools",
		},
	},
	"setfacl": {
		CheckCmd: []string{"which", "setfacl"},
		Packages: map[distroFamily]string{
			distroFedora: "acl",
			distroDebian: "acl",
			distroArch:   "acl",
			distroSUSE:   "acl",
		},
	},
	"yad": {
		CheckCmd: []string{"which", "yad"},
		Packages: map[distroFamily]string{
			distroFedora: "yad",
			distroDebian: "yad",
			distroArch:   "yad",
			distroSUSE:   "yad",
		},
	},
	"xdg-utils": {
		CheckCmd: []string{"which", "xdg-open"},
		Packages: map[distroFamily]string{
			distroFedora: "xdg-utils",
			distroDebian: "xdg-utils",
			distroArch:   "xdg-utils",
			distroSUSE:   "xdg-utils",
		},
	},
}

type distroInfo struct {
	Name       string
	Keywords   []string
	InstallCmd []string
	RemoveCmd  []string
}

var distroInfos = map[distroFamily]distroInfo{
	distroFedora: {
		Name:       "fedora",
		Keywords:   []string{"fedora", "rhel", "centos", "almalinux", "rocky"},
		InstallCmd: []string{"dnf", "install", "-y"},
		RemoveCmd:  []string{"dnf", "remove", "-y"},
	},
	distroDebian: {
		Name:       "debian",
		Keywords:   []string{"debian", "ubuntu", "kali", "mint", "pop"},
		InstallCmd: []string{"apt-get", "install", "-y"},
		RemoveCmd:  []string{"apt-get", "remove", "-y"},
	},
	distroArch: {
		Name:       "arch",
		Keywords:   []string{"arch", "manjaro", "endeavour"},
		InstallCmd: []string{"pacman", "-S", "--noconfirm"},
		RemoveCmd:  []string{"pacman", "-R", "--noconfirm"},
	},
	distroSUSE: {
		Name:       "suse",
		Keywords:   []string{"suse", "opensuse"},
		InstallCmd: []string{"zypper", "install", "-y"},
		RemoveCmd:  []string{"zypper", "remove", "-y"},
	},
}

type installedDeps struct {
	Packages []string `json:"packages"`
	Distro   string   `json:"distro"`
}

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

	for family, info := range distroInfos {
		if containsAny(combined, info.Keywords...) {
			return family
		}
	}

	return distroUnknown
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
	if info, ok := distroInfos[distro]; ok {
		return info.Name
	}
	return "unknown"
}

func distroFromName(name string) distroFamily {
	for family, info := range distroInfos {
		if info.Name == name {
			return family
		}
	}
	return distroUnknown
}

func containsAny(s string, needles ...string) bool {
	for _, needle := range needles {
		if strings.Contains(s, needle) {
			return true
		}
	}
	return false
}

func checkDeps() []string {
	headless := system.IsHeadless()

	var missing []string
	for name, dep := range deps {
		if headless && isGUIDep(name) {
			continue
		}
		if _, err := command.Output(dep.CheckCmd[0], dep.CheckCmd[1:]...); err != nil {
			missing = append(missing, name)
		}
	}
	return missing
}

func isGUIDep(name string) bool {
	switch name {
	case "yad", "xdg-utils":
		return true
	}
	return false
}

func installDeps(distro distroFamily, missing []string) error {
	info, ok := distroInfos[distro]
	if !ok {
		return fmt.Errorf("unsupported distro — install manually: %s", strings.Join(missing, ", "))
	}

	var packages []string
	for _, name := range missing {
		dep, ok := deps[name]
		if !ok {
			continue
		}
		pkg, ok := dep.Packages[distro]
		if !ok {
			continue
		}
		packages = append(packages, pkg)
	}

	if len(packages) == 0 {
		return nil
	}

	cmd := append(info.InstallCmd, packages...)

	if distro == distroDebian {
		_ = command.Run("apt-get", "update", "-qq")
	}

	if err := command.Run(cmd[0], cmd[1:]...); err != nil {
		return err
	}

	saveInstalledDeps(distro, packages)
	return nil
}

func ensureNotifierDeps() {
	missing := checkDeps()
	if len(missing) == 0 {
		return
	}

	logger.Status("installing missing notifier dependencies: " + strings.Join(missing, ", "))

	distro := detectDistro()
	if distro == distroUnknown {
		logger.Warn("could not detect distro — install manually: " + strings.Join(missing, ", "))
		return
	}

	if err := installDeps(distro, missing); err != nil {
		logger.Warn("failed to install notifier dependencies: " + err.Error())
	}
}

func removeInstalledDeps() {
	data, err := os.ReadFile(model.InstalledDepsPath)
	if err != nil {
		return
	}

	var installed installedDeps
	if err := json.Unmarshal(data, &installed); err != nil {
		return
	}

	if len(installed.Packages) == 0 {
		return
	}

	distro := distroFromName(installed.Distro)
	info, ok := distroInfos[distro]
	if !ok {
		logger.Warn("could not determine distro for dep removal")
		return
	}

	logger.Status("removing dependencies installed by statehound: " + strings.Join(installed.Packages, ", "))

	cmd := append(info.RemoveCmd, installed.Packages...)
	_ = command.Run(cmd[0], cmd[1:]...)
	_ = os.Remove(model.InstalledDepsPath)
}
