package osx

import "os"

func Getenv(key, defaultValue string) string {
	var v string

	if v = os.Getenv(key); v == "" {
		return defaultValue
	}

	return v
}
