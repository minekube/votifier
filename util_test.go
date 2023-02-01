package votifier

import (
	"strconv"
	"testing"
	"time"
)

func TestParseTime(t *testing.T) {
	now := time.Now()

	t.Run("Valid Unix Millis", func(t *testing.T) {
		unixMillis := strconv.FormatInt(now.UnixMilli(), 10)
		parsedTime := parseTime(unixMillis)
		expectedTime := now.UnixMilli()
		actualTime := parsedTime.UnixMilli()

		if expectedTime != actualTime {
			t.Errorf("Expected %v, but got %v", expectedTime, actualTime)
		}
	})

	t.Run("Invalid Unix Millis", func(t *testing.T) {
		unixMillis := "invalid"
		parsedTime := parseTime(unixMillis)
		expectedTime := now.UnixMilli()
		actualTime := parsedTime.UnixMilli()

		if expectedTime != actualTime {
			t.Errorf("Expected %v, but got %v", expectedTime, actualTime)
		}
	})

	t.Run("Older than 1 hour", func(t *testing.T) {
		oneHourAgo := now.Add(-2 * time.Hour)
		unixMillis := strconv.FormatInt(oneHourAgo.UnixMilli(), 10)
		parsedTime := parseTime(unixMillis)
		expectedTime := now.UnixMilli()
		actualTime := parsedTime.UnixMilli()

		if expectedTime != actualTime {
			t.Errorf("Expected %v, but got %v", expectedTime, actualTime)
		}
	})
}
