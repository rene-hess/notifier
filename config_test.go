package main

import (
	"strings"
	"testing"
	"time"
)

func TestParseConfigRejectsUnknownFields(t *testing.T) {
	now := time.Date(2026, time.May, 3, 9, 0, 0, 0, time.UTC)
	input := strings.NewReader(`
urgncy: normal
events:
  - time: "5m"
    message: "Team meeting"
`)

	_, err := parseConfig(now, input)
	if err == nil {
		t.Fatal("expected unknown YAML field to be rejected")
	}

	if !strings.Contains(err.Error(), "urgncy") {
		t.Fatalf("expected error to mention unknown field, got %v", err)
	}
}

func TestParseTimeStringRejectsNegativeDurations(t *testing.T) {
	now := time.Date(2026, time.May, 3, 9, 0, 0, 0, time.UTC)

	_, err := parseTimeString(now, "-5m")
	if err == nil {
		t.Fatal("expected negative duration to be rejected")
	}

	if !strings.Contains(err.Error(), "negative durations") {
		t.Fatalf("expected negative duration error, got %v", err)
	}
}

func TestParseTimeStringParsesAbsoluteTime(t *testing.T) {
	now := time.Date(2026, time.May, 3, 9, 15, 0, 0, time.UTC)

	got, err := parseTimeString(now, "12:30")
	if err != nil {
		t.Fatalf("expected absolute time to parse, got %v", err)
	}

	want := time.Date(2026, time.May, 3, 12, 30, 0, 0, time.UTC)
	if !got.Equal(want) {
		t.Fatalf("expected %v, got %v", want, got)
	}
}

func TestParseTimeStringParsesPositiveDuration(t *testing.T) {
	now := time.Date(2026, time.May, 3, 9, 15, 0, 0, time.UTC)

	got, err := parseTimeString(now, "5m")
	if err != nil {
		t.Fatalf("expected duration to parse, got %v", err)
	}

	want := now.Add(5 * time.Minute)
	if !got.Equal(want) {
		t.Fatalf("expected %v, got %v", want, got)
	}
}
