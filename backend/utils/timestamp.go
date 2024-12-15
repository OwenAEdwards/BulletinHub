package utils

import (
	"time"
)

// GetTimestamp returns the current time in a formatted string (e.g., "2024-12-14 15:04:05").
func GetTimestamp() string {
	return time.Now().Format("2006-01-02 15:04:05")
}
