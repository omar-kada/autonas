package server

import (
	"net/http"
	"os"
)

// fronendFileSystem is a custom filesystem that prevents directory listings.
type frontendFileSystem struct {
	fs http.FileSystem
}

// Open checks if the path is a directory and if it contains an index.html file.
func (nfs frontendFileSystem) Open(path string) (http.File, error) {
	f, err := nfs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	s, err := f.Stat()
	if err != nil {
		return nil, err
	}

	// If it's a directory, check for an index.html file
	if s.IsDir() {
		// Check if the directory contains an index.html
		index := path + "/index.html"
		if _, err := nfs.fs.Open(index); err != nil {
			// Close the directory file descriptor to avoid leaks
			f.Close()
			// Return a "file not found" error, which becomes a 404
			return nil, os.ErrNotExist
		}
	}
	return f, nil
}
