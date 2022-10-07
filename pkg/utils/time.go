package utils

import "time"

func GetMillis() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func TimeToMillis(t *time.Time) int64 {
	return t.UnixNano() / int64(time.Millisecond)
}
