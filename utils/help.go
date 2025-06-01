package utils

import (
	"bufio"
	"fmt"
	"path/filepath"
	"strings"
)

// GetAbsolutePath converts a relative path to an absolute path.
func GetAbsolutePath(relativePath string) (string, error) {
	absPath, err := filepath.Abs(relativePath)
	if err != nil {
		return "", err
	}
	return absPath, nil
}

func GetChoice(reader *bufio.Reader, prompt string) string {
	fmt.Print(prompt)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}
