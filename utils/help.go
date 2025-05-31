package utils

import "path/filepath"

// GetAbsolutePath converts a relative path to an absolute path.
func GetAbsolutePath(relativePath string) (string, error) {
	absPath, err := filepath.Abs(relativePath)
	if err != nil {
		return "", err
	}
	return absPath, nil
}
