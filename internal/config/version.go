package config

import "fmt"

const INTERNAL_VERSION = 100

func Version() int {
	return INTERNAL_VERSION
}

func VersionString() string {
	return fmt.Sprintf("puush r%d", INTERNAL_VERSION)
}
