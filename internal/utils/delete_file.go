package utils

import "os"

func DeleteFiles(paths ...string) {
	for _, path := range paths {
		if path != "" {
			_ = os.Remove(path) // Ignore error silently
		}
	}
}
