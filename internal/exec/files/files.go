// Package files contains utility functions for file operations.
package files

import (
	"fmt"
	"os"

	copydir "github.com/otiai10/copy"
)

// CopyServicesToPath copies all files from srcFolder to servicesPath.
func CopyServicesToPath(srcFolder, servicesPath string) (string, error) {
	if servicesPath == "" {
		return "", fmt.Errorf("SERVICES_PATH not set in config. Aborting copy")
	}
	err := copydir.Copy(srcFolder, servicesPath)
	if err != nil {
		return "", fmt.Errorf("error copying services: %v", err)
	}
	fmt.Printf("Copied all files from %s to %s\n", srcFolder, servicesPath)
	return servicesPath, nil
}

// WriteEnvFile writes the content map as key=value pairs to a .env file at filePath.
func WriteEnvFile(filePath string, content map[string]interface{}) error {
	f, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating .env file at %s: %v", filePath, err)
	}
	defer f.Close()
	for k, v := range content {
		_, err := fmt.Fprintf(f, "%s=%v\n", k, v)
		if err != nil {
			return fmt.Errorf("error writing to .env file at %s: %v", filePath, err)
		}
	}
	return nil
}
