package env

import "os"

func Get(key, value string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}

	return value
}
