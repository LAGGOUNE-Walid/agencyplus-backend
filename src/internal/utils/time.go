package utils

import "time"

func InTimeInterval(start, end, check time.Time) bool {
	if start.After(end) {
		start, end = end, start
	}
	return (check.Equal(start) || check.After(start)) && (check.Equal(end) || check.Before(end))
}
