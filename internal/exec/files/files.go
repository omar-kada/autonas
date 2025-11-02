// Package files contains utility functions for file operations.
package files

import (
	"fmt"
	"os"

	copydir "github.com/otiai10/copy"
)

// Manager defines methods for file operations.
type Manager interface {
	CopyToPath(srcFolder, servicesPath string) error
	WriteToFile(filePath string, content string) error
}

// NewManager creates a new instance of the File Manager.
func NewManager() Manager {
	return &defaultFileManager{}
}

type defaultFileManager struct{}

// CopyToPath copies all files from srcFolder to servicesPath.
func (d *defaultFileManager) CopyToPath(srcFolder, servicesPath string) error {
	if servicesPath == "" {
		return fmt.Errorf("SERVICES_PATH not set in config. Aborting copy")
	}
	err := copydir.Copy(srcFolder, servicesPath)
	if err != nil {
		return fmt.Errorf("error copying services: %v", err)
	}
	fmt.Printf("Copied all files from %s to %s\n", srcFolder, servicesPath)
	return nil
}

// WriteToFile writes the content map as key=value pairs to a .env file at filePath.
func (d *defaultFileManager) WriteToFile(filePath string, content string) error {
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
