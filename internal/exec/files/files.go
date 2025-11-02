// Package files contains utility functions for file operations.
package files

import (
	"fmt"
	"os"
)

// WriteToFile writes the content map as key=value pairs to a .env file at filePath.
func WriteToFile(filePath string, content string) error {
	f, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating .env file at %s: %v", filePath, err)
	}
	defer f.Close()
	_, err = fmt.Fprint(f, content)
	if err != nil {
		return fmt.Errorf("error writing to .env file at %s: %v", filePath, err)
	}
	return nil
}
