package helpers

import (
	"fmt"
	"strings"
	"time"
)

func ToSnakeCase(str string) string {
	lowerCasedStr := strings.ToLower(str)
	fragments := strings.Fields(lowerCasedStr)
	return strings.Join(fragments, "_")
}

func TimeToSeconds(timeStr string) (int, error) {
	colonCount := strings.Count(timeStr, ":")

	var err error
	var t time.Time

	switch colonCount {
	case 1:
		t, err = time.Parse("15:04:05", "0:"+timeStr)
		if err != nil {
			return -1, err
		}
	case 2:
		t, err = time.Parse("15:04:05", timeStr)
		if err != nil {
			return -1, err
		}
	default:
		return -1, fmt.Errorf("%s is not parsable to duration", timeStr)
	}

	return int(t.Second() + t.Minute()*60 + t.Hour()*3600), nil
}
