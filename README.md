# Schedule Notifier

A simple command-line tool that reads a YAML schedule file and sends desktop notifications at
specified times using `notify-send`.

## Requirements

- Linux system with `notify-send` installed (usually part of `libnotify-bin` package)
- A recent enough version of `go`

## Installation

```bash
git clone https://github.com/rene-hess/notifier.git
go run ./... -config example.yaml
```

## Usage

1. Create a YAML schedule file (e.g., `schedule.yaml`) with your notifications:

```yaml
icon: "/path/to/default/icon.png"  # optional default icon
urgency: normal                    # optional default urgency (low|normal|critical)

events:
  - time: "12:00"                  # Absolute time in format HH:MM
    message: "Team meeting"
    icon: "/path/to/meeting.png"   # optional override
    urgency: critical              # optional override
  - time: "5m"                     # Relative time (must parse as time.duration)
    message: "Lunch break"
```

2. Run the program:

```bash
schedule-notifier -config schedule.yaml
```

The program will run until all notifications are sent or until interrupted with Ctrl+C.

## Features

- Configurable urgency levels (low, normal, critical)
- Absolute and relative times
- Support for custom icons per notification

Reason for not having a date: This is not a calendar for reminding you of important appointments.
The main motivation was to have simple system notifications for things like switching from sitting
to standing desk position. For this purpose it makes more sense to have reusable configs where you
do not need to edit the date.


