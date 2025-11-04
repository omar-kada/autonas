// Package files contains utility functions for file operations.
package files

import (
	"fmt"
	"os"
)

// WriteToFile writes the content map as key=value pairs to a file at filePath.
func WriteToFile(filePath string, content string) error {
	f, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating file at %s: %v", filePath, err)
	}
	defer f.Close()
	_, err = fmt.Fprint(f, content)
	if err != nil {
		return fmt.Errorf("error writing to file at %s: %v", filePath, err)
	}
	return nil
}
