// Package files contains utility functions for file operations.
package files

import (
	"fmt"
	"os"
)

// Writer abstracts writing content to a file
type Writer interface {
	WriteToFile(filePath string, content string) error
}

type fileWriter struct{}

// NewWriter creates and new Writer and returns it
func NewWriter() Writer {
	return fileWriter{}
}

// WriteToFile writes the content map as key=value pairs to a file at filePath.
func (fileWriter) WriteToFile(filePath string, content string) error {
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
