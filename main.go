package main

import (
	"errors"
	"flag"
	"log"
	"log/slog"
)

func main() {
	err := run()
	if err != nil {
		log.Fatalf("Error: %v", err)
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

	slog.Info("Loaded config", "config", config)

	return nil
}
