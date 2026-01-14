package server

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestFrontendFileSystem(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()

	// Create some test files and directories
	testFiles := []struct {
		path    string
		content string
	}{
		{path: "index.html", content: "<html>Hello, World!</html>"},
		{path: "subdir/index.html", content: "<html>Subdir</html>"},
		{path: "subdir_no_index/noindex.txt", content: "No index here"},
	}

	for _, tf := range testFiles {
		fullPath := tempDir + "/" + tf.path
		if err := os.MkdirAll(fullPath[:len(fullPath)-len(tf.path[len(tf.path)-len("index.html"):])], 0740); err != nil {
			t.Fatalf("Failed to create directory: %v", err)
		}
		if err := os.WriteFile(fullPath, []byte(tf.content), 0644); err != nil {
			t.Fatalf("Failed to write file: %v", err)
		}
	}

	// Create a frontendFileSystem instance
	fs := frontendFileSystem{fs: http.Dir(tempDir)}

	// Test cases
	testCases := []struct {
		name           string
		path           string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Existing file",
			path:           "/",
			expectedStatus: http.StatusOK,
			expectedBody:   "<html>Hello, World!</html>",
		},
		{
			name:           "Directory with index.html",
			path:           "/subdir/",
			expectedStatus: http.StatusOK,
			expectedBody:   "<html>Subdir</html>",
		},
		{
			name:           "Directory without index.html",
			path:           "/subdir_no_index/",
			expectedStatus: http.StatusNotFound,
			expectedBody:   "404 page not found\n",
		},
		{
			name:           "Non-existent file",
			path:           "/nonexistent.html",
			expectedStatus: http.StatusNotFound,
			expectedBody:   "404 page not found\n",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tc.path, http.NoBody)
			rr := httptest.NewRecorder()

			// Create a handler that uses our frontendFileSystem
			handler := http.FileServer(fs)
			handler.ServeHTTP(rr, req)

			// Check the status code
			if status := rr.Code; status != tc.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tc.expectedStatus)
			}

			// Check the response body
			if rr.Body.String() != tc.expectedBody {
				t.Errorf("handler returned unexpected body: got %v want %v",
					rr.Body.String(), tc.expectedBody)
			}
		})
	}
}
