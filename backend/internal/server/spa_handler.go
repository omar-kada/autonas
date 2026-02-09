package server

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func spaHandler(w http.ResponseWriter, r *http.Request) {
	frontDir, _ := filepath.Abs(filepath.Join("frontend", "dist"))
	path := filepath.Join(frontDir, r.URL.Path)

	// check for path traversing
	absPath, err := filepath.Abs(path)
	if err != nil || !strings.HasPrefix(absPath, frontDir) {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}
	_, err = os.Stat(absPath)
	if os.IsNotExist(err) {
		http.ServeFile(w, r, filepath.Join(frontDir, "index.html"))
		return
	}

	http.FileServer(http.Dir(frontDir)).ServeHTTP(w, r)
}
