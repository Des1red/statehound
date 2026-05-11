# Statehound

Statehound is a defensive host-monitoring daemon for Linux. It runs as a root-owned systemd service, watches system state in the background, and records changes as security events.

The project is currently focused on raw visibility first: collect what changed, log it clearly, and filter/classify later.

## Current Goal

Statehound monitors:

- systemd service/unit state changes
- listening TCP/UDP ports
- process details for listening ports
- live daemon health and snapshot status

It is designed to be lightweight by using periodic snapshots instead of constant aggressive scanning.

## Architecture

```text
statehound
├── cmd
│   └── CLI flags and command routing
├── bootstrap
│   └── install, uninstall, purge, verification
├── runtime
│   └── CLI actions such as start, stop, status, events, hunt
├── daemon
│   └── root daemon entrypoint and Unix socket server
├── client
│   └── Unix socket client used by CLI commands
├── statehound
│   └── manager, snapshots, collectors, diff logic, events, hunt
├── logger
│   └── centralized output helpers
├── command
│   └── command execution helpers
└── model
    └── constants, flags, paths
```

## Runtime Model

Statehound uses one binary in two modes:

```text
statehound --daemon
    internal mode used by systemd
    owns the manager, snapshots, event generation, and socket server

statehound / shound CLI flags
    short-lived controller commands
    talk to the daemon through a Unix socket when live state is needed
```

Systemd starts the daemon:

```text
/usr/local/bin/statehound --daemon
```

The CLI talks to the daemon through:

```text
/run/statehound/statehound.sock
```

The socket is root-only:

```text
root:root 0600
```

The daemon also checks Unix peer credentials and only accepts root clients.

## Installed Paths

```text
/usr/local/bin/statehound
/usr/local/bin/shound
/etc/systemd/system/statehound.service
/etc/statehound
/var/log/statehound/events.log
/run/statehound/statehound.sock
```

`statehound` is the real binary.

`shound` is a convenience symlink alias.

## Commands

### Install

```bash
sudo go run . --install
```

Installs the binary, creates the `shound` alias, writes the systemd service, enables the service, and starts the daemon.

### Uninstall

```bash
sudo shound --uninstall
```

Removes the binary, alias, and systemd service. It keeps config and event logs.

### Purge

```bash
sudo shound --uninstall --purge
```

Removes the binary, alias, service file, config directory, and event logs.

### Start

```bash
sudo shound --start
```

Starts the systemd service if the daemon is not already responding.

### Stop

```bash
sudo shound --stop
```

Stops the systemd service and verifies the daemon socket is no longer responding.

### Restart

```bash
sudo shound --restart
```

Restarts the systemd service and verifies the daemon responds after restart.

### Status

```bash
sudo shound --status
```

Queries the daemon through the Unix socket and prints live manager state.

Example:

```text
[+] statehound daemon is running
manager=running
interval=5s
last_scan=2026-05-11T02:24:59+03:00
systemd_services=217
active_services=87
listening_ports=26
```

### System Logs

```bash
sudo shound --logs
```

Shows systemd journal logs for the Statehound service.

### Events

```bash
sudo shound --events
```

Shows Statehound detection events from:

```text
/var/log/statehound/events.log
```

### Clear Events

```bash
sudo shound --clear-events
```

Clears the Statehound event log.

### Hunt

```bash
sudo shound --hunt <target>
```

Searches the daemon's live snapshot for matching services and listening ports.

Examples:

```bash
sudo shound --hunt proxychan
sudo shound --hunt 8000
sudo shound --hunt nc
sudo shound --hunt sshd.service
```

Example output:

```text
[+] hunt target: proxychan

Services:
  proxychan.service desc="ProxyChan SOCKS5 Proxy" load=loaded state=active/running

Listening ports:
  tcp 127.0.0.1:1080 scope=local process=proxychan pid=83184 exe=/usr/local/bin/proxychan cmdline="/usr/local/bin/proxychan --listen 127.0.0.1:1080 ..."
  tcp 127.0.0.1:6060 scope=local process=proxychan pid=83184 exe=/usr/local/bin/proxychan cmdline="/usr/local/bin/proxychan --listen 127.0.0.1:1080 ..."
```

## Events

Statehound records events in a simple line-based format:

```text
[time] TYPE message
```

Example:

```text
[2026-05-11T03:01:14+03:00] SERVICE_STARTED proxychan.service desc="ProxyChan SOCKS5 Proxy" load=loaded state=active/running transition=inactive/dead->active/running
[2026-05-11T03:01:19+03:00] PORT_OPENED tcp 127.0.0.1:1080 scope=local process=proxychan pid=84761 exe=/usr/local/bin/proxychan cmdline="/usr/local/bin/proxychan --listen 127.0.0.1:1080 ..."
```

## Event Types

### Manager Events

```text
MANAGER_STARTED
BASELINE_CREATED
```

### Systemd Unit Events

```text
UNIT_APPEARED
UNIT_REMOVED
```

These mean the unit appeared or disappeared from systemd's loaded-unit list.

### Service State Events

```text
SERVICE_STARTED
SERVICE_STOPPED
SERVICE_FAILED
SERVICE_STATE_CHANGED
```

These mean a known unit changed runtime state.

### Port Events

```text
PORT_OPENED
PORT_CLOSED
```

These are generated when a listening TCP/UDP socket appears or disappears.

## Service Monitoring

Statehound collects all loaded systemd service units using:

```bash
systemctl list-units --type=service --all --no-pager --no-legend
```

It tracks:

```text
name
load state
active state
sub state
description
```

This allows Statehound to detect more than only running services. It can see transitions such as:

```text
inactive/dead -> active/running
active/running -> inactive/dead
active/running -> failed/failed
activating/start -> active/exited
```

The status output includes both total systemd service units and active service count.

## Port Monitoring

Statehound collects listening TCP/UDP sockets using `ss`.

For each listening port, it records:

```text
protocol
address
port
scope
process name
PID
executable path
command line
```

Port scope is classified as:

```text
local      127.0.0.1, ::1
public     0.0.0.0, *, ::
interface  specific interface IPs
```

Examples:

```text
tcp 127.0.0.1:9001 scope=local process=idor_server pid=34446
tcp 0.0.0.0:4444 scope=public process=nc pid=48599
udp 10.240.175.43:3702 scope=interface process=wsdd pid=57631
```

Long command lines are truncated to keep events readable.

## Security Model

Statehound is intended to be root-controlled.

The daemon runs as root through systemd:

```ini
User=root
Group=root
```

Important files and directories are root-owned:

```text
/usr/local/bin/statehound      root:root 0755
/usr/local/bin/shound          root:root symlink
/etc/systemd/system/...        root:root 0644
/etc/statehound                root:root 0700
/var/log/statehound            root:root 0700
/run/statehound                root:root 0700
/run/statehound/statehound.sock root:root 0600
```

The Unix socket also verifies peer credentials and rejects non-root clients.

This means normal users cannot query, control, or manipulate Statehound runtime state.

## Current Stage

Statehound currently has:

```text
Stage 1:
- install/uninstall/purge
- systemd service lifecycle
- root-owned daemon
- Unix socket communication
- root-only socket access
- status/logs/events commands
- shound alias

Stage 2:
- manager loop
- baseline snapshots
- systemd service state monitoring
- listening port monitoring
- service/unit diff events
- port open/close diff events
- process enrichment for listening ports
- port scope classification
- live status details
- hunt mode
```

## Planned Stage 3

Stage 3 will focus on filtering and classification.

The current behavior intentionally logs raw truth. This creates noise from normal desktop activity such as browser UDP sockets, user session units, NetworkManager dispatcher events, and discovery services.

Planned Stage 3 ideas:

```text
- event severity levels
- noise classification
- suspicious listener classification
- ignore rules
- raw vs filtered event views
- default event tailing
- rules for browser UDP noise
- rules for user@*.service noise
- rules for NetworkManager-dispatcher noise
- high severity for public nc/ncat listeners
- high severity for listeners from /tmp or /dev/shm
```

A likely future model:

```text
raw event created
classifier tags event
filter decides display/storage behavior
filtered events shown by default
raw events available on demand
```

## Development Notes

Statehound is designed around simple responsibility boundaries:

```text
collectors collect state
manager owns snapshots
diff functions detect changes
event writer persists events
socket exposes live daemon state
runtime package implements CLI actions
```

The current design intentionally avoids full process watching because that would create heavy noise. Instead, process enrichment is attached to listening ports, where it is most useful defensively.