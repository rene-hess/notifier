package main

import (
	"context"
	"errors"
	"flag"
	"log/slog"
	"os"
	"os/exec"
	"os/signal"
	"sort"
	"syscall"
	"time"
)

func main() {
	err := run()
	if err != nil {
		slog.Error("Error running notifier", "error", err)
		os.Exit(1)
	}
}

func run() error {
	yamlPath := flag.String(
		"config",
		"schedule.yaml",
		"Path to the YAML schedule file",
	)
	flag.Parse()

	if *yamlPath == "" {
		return errors.New("config path is required, call with --config <path to config yaml>")
	}

	config, err := loadConfig(*yamlPath)
	if err != nil {
		return err
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer func() {
		slog.Info("Stoping notifier")
		stop()
	}()

	slog.Info("Starting notifier")
	slog.Debug("With config", "config", config)
	notify(ctx, config)

	return nil
}

func notify(ctx context.Context, config Config) {
	config.Events = sortEvents(config.Events)
	for _, event := range config.Events {
		wait := time.Until(event.Time)
		if wait > 0 {
			timer := time.NewTimer(wait)
			select {
			case <-ctx.Done():
				timer.Stop()
				return
			case <-timer.C:
			}
		}
		args := []string{event.Message}
		if event.Urgency != nil {
			args = append([]string{"-u", string(*event.Urgency)}, args...)
		} else {
			args = append([]string{"-u", string(config.Urgency)}, args...)
		}
		cmd := exec.Command("notify-send", args...)
		err := cmd.Run()
		if err != nil {
			slog.Error("Error running notify-send", "error", err)
		}
	}
}

func sortEvents(events []Event) []Event {
	sort.Slice(events, func(i, j int) bool {
		return events[i].Time.Before(events[j].Time)
	})
	return events
}
