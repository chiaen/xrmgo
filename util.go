package xrmgo

import (
	"time"
)

func toCurrentTime(t time.Time) string {
	return t.UTC().Format(time.RFC3339)
}

func toTomorrowTime(t time.Time) string {
	return t.UTC().Add(time.Hour * 24).Format(time.RFC3339)
}
