package utils

import "time"

// GetTimeStamp возвращает unix ts
func GetTimeStamp(t time.Time) int64 {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC).Unix()
}
