package votifier

import (
	"testing"
	"time"
)

func TestParseTime(t *testing.T) {
	now := time.Now()
	timeNow = func() time.Time { return now }

	t.Run("valid", func(t *testing.T) {
		unixMillis := "1609459200000"
		expected := now
		result := parseTime(unixMillis)
		if !result.Equal(expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("invalid", func(t *testing.T) {
		unixMillis := "" // invalid
		expected := now
		result := parseTime(unixMillis)
		if !result.Equal(expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("older than 1 hour", func(t *testing.T) {
		unixMillis := "0"
		expected := now
		result := parseTime(unixMillis)
		if !result.Equal(expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})
}
