package utils

import "time"

func GetCurrentTimeInUnix() int64 {
	return time.Now().Unix()
}
