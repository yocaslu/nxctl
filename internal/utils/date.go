package utils

import "time"

func GetDate() string {
	return time.Now().Format("02-01-2006")
}
