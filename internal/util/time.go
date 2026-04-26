package util

import "time"

func NowNano() int64 {
	return time.Now().UnixNano()
}
