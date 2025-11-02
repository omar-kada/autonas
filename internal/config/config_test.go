package config

import (
	"embed"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"golang.org/x/tools/txtar"
	"gopkg.in/yaml.v3"
)

//go:embed test_data/*
var testDataFS embed.FS

func TestConfigPerService_BuildsCorrectMap(t *testing.T) {
	cfg := Config{
		AutonasHost:  "host",
		ServicesPath: "/services",
		DataPath:     "/data",
		Extra:        map[string]any{"GLOBAL": "g"},
		Services: map[string]ServiceConfig{
			"svc": {
				Port:    8080,
				Version: "v1",
				Extra:   map[string]any{"SVC_EXTRA": "s"},
			},
		},
	}

	got := cfg.PerService("svc")
	want := map[string]any{
		"AUTONAS_HOST":  "host",
		"SERVICES_PATH": "/services",
		"DATA_PATH":     "/data/svc",
		"GLOBAL":        "g",
		"PORT":          8080,
		"VERSION":       "v1",
		"SVC_EXTRA":     "s",
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("ConfigPerService mismatch\nwant=%#v\ngot =%#v", want, got)
	}
}

// getTempTestFile writes an embedded fixture (from testDataFS) to a temp file and returns its path.
func getTempTestFile(t *testing.T, name string) string {
	t.Helper()
	bs, err := testDataFS.ReadFile("test_data/" + name)
	if err != nil {
		t.Fatalf("read fixture %s: %v", name, err)
	}
	dst := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(dst, bs, 0o644); err != nil {
		t.Fatalf("write fixture: %v", err)
	}
	return dst
}

func TestDecodeConfig(t *testing.T) {
	input := map[string]any{
		"AUTONAS_HOST": "localhost",
		"DATA_PATH":    "/data",
		"enabled_services": []any{
			"svc1"},
		"services": map[string]any{
			"svc1": map[string]any{
				"PORT":      8080,
				"VERSION":   "v1",
				"NEW_FIELD": "new_value",
			},
			"svc2": map[string]any{
				"PORT":    9090,
				"VERSION": "v2",
			},
		},
	}
	want := Config{
		AutonasHost:     "localhost",
		ServicesPath:    "",
		DataPath:        "/data",
		EnabledServices: []string{"svc1"},
		Services: map[string]ServiceConfig{
			"svc1": {
				Port:    8080,
				Version: "v1",
				Extra:   map[string]any{"NEW_FIELD": "new_value"},
			},
			"svc2": {
				Port:    9090,
				Version: "v2",
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

// extractTxtar extracts config files & expected result from a txtar test archive.
func extractTxtar(t *testing.T, archivePath string) (inputs []string, want map[string]any) {
	t.Helper()

	ar, err := txtar.ParseFile(getTempTestFile(t, archivePath))
	if err != nil {
		t.Fatalf("failed to parse txtar file: %v", err)
	}

	tempDir := t.TempDir()
	for _, file := range ar.Files {
		outPath := filepath.Join(tempDir, file.Name)
		if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
			t.Fatalf("failed to create tmp directory: %v", err)
		}
		if err := os.WriteFile(outPath, file.Data, 0o600); err != nil {
			t.Fatalf("failed to create tmp file: %v", err)
		}

		if file.Name == "want" {
			err = yaml.Unmarshal(file.Data, &want)
			if err != nil {
				t.Fatalf("failed to parse want yaml file: %v", err)
			}
		} else {
			inputs = append(inputs, outPath)
		}
	}
	return inputs, want
}

func TestLoadConfig_SuccessWithOverride(t *testing.T) {
	testTable := []string{"config_override.txtar"}

	for _, testFile := range testTable {
		t.Run(testFile, func(t *testing.T) {

			inputs, want := extractTxtar(t, testFile)
			wantCfg, err := decodeConfig(want)
			if err != nil {
				t.Fatalf("failed to decode config: %v", err)
			}
			cfg, err := FromFiles(inputs)
			if err != nil {
				t.Fatalf("LoadConfig failed: %v", err)
			}
			if !reflect.DeepEqual(cfg, wantCfg) {
				t.Fatalf("unexpected svc values: want=%+v got=%+v", wantCfg, cfg)
			}
		})
	}

}

func TestLoadConfig_FileError(t *testing.T) {
	t.Run("missing file", func(t *testing.T) {
		if _, err := FromFiles([]string{"/does/not/exist.yaml"}); err == nil {
			t.Fatalf("expected error for missing file")
		}
	})

	t.Run("invalid yaml", func(t *testing.T) {
		// create a temporary file with invalid YAML content so the test
		tmp := t.TempDir()
		f := filepath.Join(tmp, "invalid.yaml")
		invalid := []byte("this: is: not: valid: yaml")
		if err := os.WriteFile(f, invalid, 0o644); err != nil {
			t.Fatalf("failed to write temp invalid yaml: %v", err)
		}
		if _, err := FromFiles([]string{f}); err == nil {
			t.Fatalf("expected error for invalid yaml")
		}
	})
}
