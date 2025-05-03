package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

type Urgency string

const (
	low      Urgency = "low"
	normal   Urgency = "normal"
	critical Urgency = "critical"
)

type Event struct {
	TimeString string  `yaml:"time"`
	Message    string  `yaml:"message"`
	Urgency    Urgency `yaml:"urgency,omitempty"`
	Icon       string  `yaml:"icon,omitempty"`

	// Internal field for parsed time, not mapped to YAML
	Time time.Time `yaml:"-"`
}

type Config struct {
	// Default values for the events
	Urgency Urgency `yaml:"urgency,omitempty"`
	Icon    string  `yaml:"icon,omitempty"`

	Events []Event `yaml:"events"`
}

func loadConfig(path string) (Config, error) {
	now := time.Now()

	file, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("error reading yaml file %s: %w", path, err)
	}

	config, err := parseConfig(now, bytes.NewReader(file))
	if err != nil {
		return Config{}, err
	}

	return config, nil
}

func parseConfig(now time.Time, input io.Reader) (Config, error) {
	var config Config
	decoder := yaml.NewDecoder(input)
	if err := decoder.Decode(&config); err != nil {
		return Config{}, fmt.Errorf("failed to decode yaml: %w", err)
	}

	// Set default urgency if not specified
	if config.Urgency == "" {
		config.Urgency = normal
	}

	err := validateConfig(config)
	if err != nil {
		return Config{}, fmt.Errorf("invalid config: %w", err)
	}

	for i, event := range config.Events {
		targetTime, err := parseTimeString(now, event.TimeString)
		if err != nil {
			return Config{}, fmt.Errorf("invalid time string %s: %w", event.TimeString, err)
		}
		config.Events[i].Time = targetTime
	}

	return config, nil
}

func isValidUrgency(urgency Urgency) bool {
	switch urgency {
	case low, normal, critical:
		return true
	default:
		return false
	}
}

func validateConfig(config Config) error {
	if !isValidUrgency(config.Urgency) {
		return fmt.Errorf("invalid urgency: %s", config.Urgency)
	}
	if config.Icon != "" {
		if err := pathValid(config.Icon); err != nil {
			return fmt.Errorf("icon file does not exist: %s: %w", config.Icon, err)
		}
	}

	if len(config.Events) == 0 {
		return fmt.Errorf("no events found")
	}

	for _, event := range config.Events {
		if event.TimeString == "" {
			return fmt.Errorf("event time is required")
		}
		if event.Message == "" {
			return fmt.Errorf("event message is required")
		}
		if event.Urgency != "" && !isValidUrgency(event.Urgency) {
			return fmt.Errorf("invalid urgency: %s", event.Urgency)
		}
		if event.Icon != "" {
			if err := pathValid(event.Icon); err != nil {
				return fmt.Errorf("icon file does not exist: %s: %w", event.Icon, err)
			}
		}
	}

	return nil
}

func pathValid(path string) error {
	if !filepath.IsAbs(path) {
		return fmt.Errorf("icon path must be absolute: %s", path)
	}
	_, err := os.Stat(path)
	if err != nil {
		return err
	}

	return nil
}

func parseTimeString(now time.Time, timeStr string) (time.Time, error) {
	// Try parsing as absolute time (HH:MM)
	if t, err := time.Parse("15:04", timeStr); err == nil {
		return time.Date(
			now.Year(),
			now.Month(),
			now.Day(),
			t.Hour(),
			t.Minute(),
			0, 0,
			now.Location(),
		), nil
	}

	// Try parsing relative time
	d, err := time.ParseDuration(timeStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid time string: %w", err)
	}

	return now.Add(d), nil
}
