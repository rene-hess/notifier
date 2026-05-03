package main

import (
	"reflect"
	"testing"
	"time"
)

func TestNotifyArgsTreatsMessageAsPositional(t *testing.T) {
	config := Config{
		Urgency: normal,
		Icon:    "/default/icon.png",
	}
	event := Event{
		Message: "--help",
		Urgency: critical,
		Icon:    "/event/icon.png",
	}

	got := notifyArgs(config, event)
	want := []string{"-i", "/event/icon.png", "-u", "critical", "--", "--help"}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expected %v, got %v", want, got)
	}
}

func TestSortEventsPreservesOrderForEqualTimes(t *testing.T) {
	sharedTime := time.Date(2026, time.May, 3, 16, 0, 0, 0, time.UTC)
	events := []Event{
		{Message: "second", Time: sharedTime},
		{Message: "first", Time: time.Date(2026, time.May, 3, 8, 0, 0, 0, time.UTC)},
		{Message: "third", Time: sharedTime},
	}

	sorted := sortEvents(events)
	got := []string{sorted[0].Message, sorted[1].Message, sorted[2].Message}
	want := []string{"first", "second", "third"}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expected %v, got %v", want, got)
	}
}
