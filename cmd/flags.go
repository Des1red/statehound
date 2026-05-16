package cmd

import (
	"os"

	"statehound/internal/model"

	"github.com/spf13/pflag"
)

func flags() model.Flags {
	var f model.Flags

	fs := pflag.NewFlagSet("statehound", pflag.ExitOnError)

	fs.BoolVarP(&f.CLI.Help, "help", "h", false, "show help")
	fs.BoolVarP(&f.CLI.Version, "version", "v", false, "show version")

	fs.BoolVar(&f.CLI.Install, "install", false, "install statehound binary and systemd service")
	fs.BoolVar(&f.CLI.Uninstall, "uninstall", false, "remove statehound binary and systemd service")
	fs.BoolVar(&f.CLI.Purge, "purge", false, "remove statehound config and event logs during uninstall")

	fs.BoolVar(&f.CLI.Start, "start", false, "start statehound service")
	fs.BoolVar(&f.CLI.Stop, "stop", false, "stop statehound service")
	fs.BoolVar(&f.CLI.Restart, "restart", false, "restart statehound service")
	fs.BoolVar(&f.CLI.Status, "status", false, "show statehound daemon status")
	fs.BoolVar(&f.CLI.Logs, "logs", false, "show statehound logs")
	fs.BoolVar(&f.CLI.Events, "events", false, "show statehound events")
	fs.BoolVar(&f.CLI.EventsGUI, "events-gui", false, "open events viewer")
	fs.StringVarP(&f.CLI.Filter, "filter", "f", "", "filter events by urgency: critical, normal")
	fs.BoolVar(&f.CLI.ClearEvents, "clear-events", false, "clear statehound event log")
	fs.BoolVar(&f.CLI.Snapshot, "snapshot", false, "show current statehound snapshot")
	fs.StringVar(&f.CLI.Hunt, "hunt", "", "investigate a service, process, pid, or port")

	fs.BoolVar(&f.Runtime.Daemon, "daemon", false, "run statehound daemon")
	fs.BoolVar(&f.Runtime.Notify, "notify", false, "run statehound desktop notifier")
	fs.BoolVar(&f.CLI.Verbose, "verbose", false, "enable verbose output")

	_ = fs.Parse(os.Args[1:])

	syncRuntimeFlags(&f)

	return f
}

func syncRuntimeFlags(f *model.Flags) {
	f.Runtime.Verbose = f.CLI.Verbose
}

func handleFlags() {
	f := flags()

	switch {
	case f.CLI.Help:
		printHelp()
	case f.CLI.Version:
		printVersion()

	case f.Runtime.Daemon:
		statehoundd()
	case f.Runtime.Notify:
		notifier()

	case f.CLI.Install:
		install()
	case f.CLI.Uninstall:
		uninstall(f.CLI.Purge)

	case f.CLI.Start:
		start()
	case f.CLI.Stop:
		stop()
	case f.CLI.Restart:
		restart()
	case f.CLI.Status:
		status()
	case f.CLI.Logs:
		logs()
	case f.CLI.Snapshot:
		snapshot()
	case f.CLI.ClearEvents:
		clearEvents()
	case f.CLI.Events:
		events(f.CLI.Filter)
	case f.CLI.EventsGUI:
		eventsGUI(f.CLI.Filter)
	case f.CLI.Hunt != "":
		hunt(f.CLI.Hunt)
	default:
		printHelp()
	}
}
