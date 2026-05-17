# Statehound

Most attacks don't announce themselves. They blend in — a new service that wasn't there yesterday, a port quietly opened by something running from `/tmp`, a cron job added while you were looking elsewhere. By the time you notice, the window has passed.

Statehound exists because Linux systems change constantly, and most of those changes go completely unobserved.

It runs as a root-owned daemon in the background, takes periodic snapshots of your system's state, and records every meaningful change as a security event. No agents. No cloud. No configuration required to get started. Just a clear, honest log of what your system is doing.

---

## What it watches

- **Systemd services** — every unit state transition: started, stopped, failed, appeared, or removed
- **Listening ports** — every TCP/UDP socket that opens or closes, with the process, PID, executable, and command line behind it
- **Persistence files** — cron jobs, systemd unit files, SSH authorized keys, shell profiles, autostart entries
- **Outbound connections** — established TCP connections and the processes behind them
- **Suspicious processes** — executables running from `/tmp`, deleted binaries, hidden home directory executables, reverse shell patterns, netcat, script servers, root processes outside standard paths

When something changes, statehound logs it. When something looks dangerous, it notifies your desktop.

---

## What it is not

Statehound is not an EDR. It does not block anything. It does not send data anywhere. It does not run kernel modules or eBPF probes.

It is a visibility tool. It tells you what changed. What you do with that is up to you.

---

## Installation

**Requirements:** Linux with systemd, Go 1.21+

```bash
git clone https://github.com/yourname/statehound
cd statehound
sudo go run . --install
```

This builds the binary, installs it to `/usr/local/bin/statehound`, creates the systemd service, enables it, and starts the daemon. A `shound` alias is also created for convenience.

Missing dependencies (yad, notify-send, acl) are installed automatically for Fedora, Debian/Ubuntu/Kali, Arch, and openSUSE.

---

## Commands

```bash
sudo shound --status               # live daemon status
sudo shound --events               # show all recorded events
sudo shound --events -f critical   # show critical events only
sudo shound --events -f normal     # show normal events only
sudo shound --events-gui           # open graphical event viewer
sudo shound --snapshot             # full current system snapshot
sudo shound --hunt <target>        # search live state by name, port, or pid
sudo shound --logs                 # systemd journal logs for the daemon
sudo shound --clear-events         # clear the event log
sudo shound --restart              # restart the daemon
sudo shound --stop                 # stop the daemon
sudo shound --start                # start the daemon
```

### Hunt

Hunt lets you investigate something specific against live state:

```bash
sudo shound --hunt nc
sudo shound --hunt 4444
sudo shound --hunt sshd.service
sudo shound --hunt proxychan
```

---

## Events

Events are written to `/var/log/statehound/events.log` in a simple, readable format:

```
[2026-05-11T03:01:14+03:00] SERVICE_STARTED [normal] proxychan.service desc="ProxyChan SOCKS5 Proxy" load=loaded state=active/running transition=inactive/dead->active/running
[2026-05-11T03:01:19+03:00] PORT_OPENED [normal] tcp 127.0.0.1:1080 scope=local process=proxychan pid=84761 exe=/usr/local/bin/proxychan
[2026-05-11T03:04:02+03:00] FILE_ADDED [critical] /etc/cron.d/backdoor size=42 mode=-rw-r--r-- hash=9f3a...
```

The log rotates automatically at 5 MB and backups are kept for 30 days.

### Urgency levels

Every event is classified as `critical` or `normal`:

**Critical** — requires immediate attention:
- Files added to cron, systemd units, SSH authorized keys, shell profiles
- Processes running from `/tmp`, `/dev/shm`, or hidden home directories
- Deleted executables still running
- Root processes outside standard system paths
- Shell tools making external connections

**Normal** — worth knowing about:
- Services starting, stopping, or failing
- Public listeners opened
- Processes running from `~/.local/bin`, script servers, network tools

Noise is filtered automatically — browser UDP traffic, known system services, and standard system activity are suppressed before logging.

---

## Desktop notifications

Statehound includes a desktop notifier that runs as your user session service. It watches for new events and sends desktop notifications grouped by urgency.

Clicking a notification opens a live graphical viewer showing both critical and normal tabs. The viewer updates automatically as new events arrive. You can also open it manually:

```bash
sudo shound --events-gui           # open on critical tab
sudo shound --events-gui -f normal # open on normal tab
```

The notifier is installed automatically for the detected desktop user when running `--install`.

---

## Uninstall

```bash
sudo shound --uninstall          # remove binary and service, keep logs
sudo shound --uninstall --purge  # remove everything including logs
```

---

## Security model

The daemon runs as root. The control socket is root-only and verifies peer credentials on every connection. Normal users cannot query or influence statehound's runtime state.

The `statehound` group is created during install. The desktop user is added to this group, giving them read access to the event log and the notification socket — nothing else. No other users can access these.
