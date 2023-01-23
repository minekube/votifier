package votifier

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net"
	"strconv"
	"time"
)

func randomString() (string, error) {
	p := make([]byte, 24)
	_, err := rand.Read(p)
	if err != nil {
		return "", err
	}

	return base64.RawStdEncoding.EncodeToString(p), nil
}

func dial(addr string) (net.Conn, error) {
	conn, err := net.DialTimeout("tcp", addr, 3*time.Second)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to %s: %w", addr, err)
	}
	err = conn.SetDeadline(timeNow().Add(3 * time.Second))
	if err != nil {
		return nil, fmt.Errorf("failed to set deadline: %v", err)
	}
	return conn, nil
}

var timeNow = time.Now

func parseTime(unixMillis string) time.Time {
	now := timeNow()
	ms, err := strconv.ParseInt(unixMillis, 10, 64)
	if err != nil {
		unix := time.UnixMilli(ms)
		// some vote sites don't sent a timestamp,
		// fallback to now if older than 1 hour
		if now.Sub(unix).Abs() < time.Hour {
			return unix
		}
	}
	return now
}

func formatTimeMillis(t time.Time) string {
	return strconv.FormatInt(t.UnixMilli(), 10)
}
