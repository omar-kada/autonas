package storage

import (
	"omar-kada/autonas/models"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecodeConfig(t *testing.T) {
	input := map[string]any{
		"AUTONAS_HOST": "localhost",
		"DATA_PATH":    "/data",
		"services": map[string]any{
			"svc1": map[string]any{
				"PORT":      8080,
				"VERSION":   "v1",
				"NEW_FIELD": "new_value",
			},
			"svc2": map[string]any{
				"Disabled": true,
				"Port":     9090,
				"Version":  "v2",
			},
		},
	}
	want := models.Config{
		Extra: map[string]any{
			"AUTONAS_HOST": "localhost",
			"DATA_PATH":    "/data",
		},
		Services: map[string]models.ServiceConfig{
			"svc1": {
				"NEW_FIELD": "new_value",
				"PORT":      8080,
				"VERSION":   "v1",
			},
			"svc2": {
				"Disabled": true,
				"Port":     9090,
				"Version":  "v2",
			},
		},
	}

	cfg, err := decodeConfig(input)
	if err != nil {
		t.Fatalf("decodeConfig failed: %v", err)
	}

	if !reflect.DeepEqual(cfg, want) {
		t.Fatalf("decodeConfig mismatch\nwant=%#v\ngot =%#v", want, cfg)
	}
}

func TestUpdateConfig(t *testing.T) {
	t.Run("successful update", func(t *testing.T) {
		tmpDir := t.TempDir()
		filePath := filepath.Join(tmpDir, "config.yaml")
		store := NewConfigStore(filePath)

		input := models.Config{
			Extra: map[string]any{
				"AUTONAS_HOST": "localhost",
				"DATA_PATH":    "/data",
			},
			Services: map[string]models.ServiceConfig{
				"svc1": {
					"PORT":    8080,
					"VERSION": "v1",
				},
			},
		}

		err := store.Update(input)
		assert.NoError(t, err)

		// Verify the file was written correctly
		_, err = os.ReadFile(filePath)
		assert.NoError(t, err)

		cfg, err := store.Get()
		assert.NoError(t, err)

		// Use deepEqual to compare the result with the expected value
		if !reflect.DeepEqual(cfg, input) {
			t.Errorf("Expected %v, got %v", input, cfg)
		}
	})

	t.Run("file write error", func(t *testing.T) {
		// Create a directory that we can't write to
		tmpDir := t.TempDir()
		readOnlyDir := filepath.Join(tmpDir, "readonly")
		err := os.Mkdir(readOnlyDir, 0555)
		assert.NoError(t, err)

		filePath := filepath.Join(readOnlyDir, "config.yaml")
		store := NewConfigStore(filePath)

		input := models.Config{
			Extra: map[string]any{
				"AUTONAS_HOST": "localhost",
			},
		}

		err = store.Update(input)
		assert.Error(t, err)
	})
}

func TestLoadConfig_FileError(t *testing.T) {
	t.Run("missing file", func(t *testing.T) {
		configStore := NewConfigStore("/does/not/exist.yaml")
		_, err := configStore.Get()
		assert.Error(t, err)
	})

	t.Run("invalid yaml", func(t *testing.T) {
		// create a temporary file with invalid YAML content so the test
		tmp := t.TempDir()
		f := filepath.Join(tmp, "invalid.yaml")
		invalid := []byte("this: is: not: valid: yaml")
		err := os.WriteFile(f, invalid, 0o644)
		assert.NoError(t, err)

		configStore := NewConfigStore(f)

		_, err = configStore.Get()
		assert.Error(t, err)
	})
}
