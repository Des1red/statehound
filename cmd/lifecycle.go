package cmd

import (
	"fmt"

	"statehound/internal/bootstrap"
	"statehound/internal/daemon"
	"statehound/internal/model"
	"statehound/internal/runtime"
)

func printHelp() {
	fmt.Println("statehound - defensive host monitoring")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  statehound --install")
	fmt.Println("  statehound --uninstall")
	fmt.Println("  statehound --uninstall --purge")
	fmt.Println("  statehound --start")
	fmt.Println("  statehound --stop")
	fmt.Println("  statehound --restart")
	fmt.Println("  statehound --status")
	fmt.Println("  statehound --logs")
	fmt.Println("  statehound --events")
	fmt.Println("  statehound --clear-events")
	fmt.Println("  statehound --hunt sshd.service")
	fmt.Println("  statehound --hunt 4444")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  -h, --help            show help")
	fmt.Println("  -v, --version         show version")
	fmt.Println("      --install         install statehound binary and systemd service")
	fmt.Println("      --uninstall       remove statehound binary and systemd service")
	fmt.Println("      --purge           remove config and event logs during uninstall")
	fmt.Println("      --start           start statehound service")
	fmt.Println("      --stop            stop statehound service")
	fmt.Println("      --restart         restart statehound service")
	fmt.Println("      --status          show statehound daemon status")
	fmt.Println("      --logs            show statehound logs")
	fmt.Println("      --events          show statehound events")
	fmt.Println("      --clear-events    clear statehound event log")
	fmt.Println("      --hunt string     investigate a service, process, pid, or port")
	fmt.Println("      --verbose         enable verbose output")
}

func printVersion() {
	fmt.Println(model.Version)
}

func install() {
	bootstrap.Install()
}

func uninstall(purge bool) {
	bootstrap.Uninstall(purge)
}

func start() {
	if !bootstrap.VerifyInstallation() {
		return
	}
	runtime.Start()
}

func stop() {
	if !bootstrap.VerifyInstallation() {
		return
	}
	runtime.Stop()
}

func restart() {
	if !bootstrap.VerifyInstallation() {
		return
	}
	runtime.Restart()
}

func status() {
	if !bootstrap.VerifyInstallation() {
		return
	}
	runtime.Status()
}

func logs() {
	if !bootstrap.VerifyInstallation() {
		return
	}
	runtime.Log()
}

func events() {
	if !bootstrap.VerifyInstallation() {
		return
	}
	runtime.Events()
}

func clearEvents() {
	if !bootstrap.VerifyInstallation() {
		return
	}
	runtime.ClearEvents()
}

func statehoundd() {
	if !bootstrap.VerifyInstallation() {
		return
	}
	daemon.Run()
}

func hunt(target string) {
	if !bootstrap.VerifyInstallation() {
		return
	}

	runtime.Hunt(target)
}

func notifier() {
	runtime.Notify()
}
