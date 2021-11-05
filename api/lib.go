package main

import (
	"os"
	"time"
)

// timestamp returns a formated timestamp
func timestamp() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

// timeDuration returns a duration object from the provided
// config value duration.
func timeDuration(duration int) time.Duration {
	if duration == 0 {
		return 0
	} else if duration == -1 {
		return -1
	}

	return time.Duration(duration) * time.Second
}

// pathEnsure creates the directory path if it
// doesn't exists.
func pathEnsure(path string) {
	if !pathExists(path) {
		os.MkdirAll(path, os.ModePerm)
	}
}

// pathExists returns true if the provided path exists.
func pathExists(path string) bool {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
