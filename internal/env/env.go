package env

import (
	"os"
	"strconv"
)
func Getstring(key string, defaultvalue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultvalue
}
func Getint(key string, defaultvalue int) int {
	if value, ok := os.LookupEnv(key); ok {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultvalue
}
