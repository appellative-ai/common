package fmtx

import (
	"strconv"
	"strings"
	"time"
)

func ParseDuration(s string) (time.Duration, error) {
	if s == "" {
		return 0, nil
	}
	tokens := strings.Split(s, "ms")
	if len(tokens) == 2 {
		val, err := strconv.Atoi(tokens[0])
		if err != nil {
			return 0, err
		}
		return time.Duration(val) * time.Millisecond, nil
	}
	tokens = strings.Split(s, "µs")
	if len(tokens) == 2 {
		val, err := strconv.Atoi(tokens[0])
		if err != nil {
			return 0, err
		}
		return time.Duration(val) * time.Microsecond, nil
	}
	tokens = strings.Split(s, "m")
	if len(tokens) == 2 {
		val, err := strconv.Atoi(tokens[0])
		if err != nil {
			return 0, err
		}
		return time.Duration(val) * time.Minute, nil
	}
	// Assume seconds
	tokens = strings.Split(s, "s")
	if len(tokens) == 2 {
		s = tokens[0]
	}
	val, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	return time.Duration(val) * time.Second, nil
}

// Milliseconds - convert time.Duration to milliseconds
func Milliseconds(duration time.Duration) int {
	if duration <= 0 {
		return -1
	}
	return int(duration / time.Duration(1e6))
}
