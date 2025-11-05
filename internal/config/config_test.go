package config

import (
	"embed"
	"omar-kada/autonas/internal/testutil"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"testing"
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

func extractTestCase(t *testing.T, archivePath string) (inputs []string, want map[string]any) {
	t.Helper()

	filesMap := testutil.ExtractTxtar(t, testDataFS, archivePath)
	for _, file := range filesMap {
		if file.Name == "want" {
			want = testutil.ReadYamlFile(t, file.Path)
		} else {
			inputs = append(inputs, file.Path)
		}
	}
	sort.Strings(inputs)
	return inputs, want
}

func TestLoadConfig_SuccessWithOverride(t *testing.T) {
	testTable := []string{"config_override.txtar"}

	for _, testFile := range testTable {
		t.Run(testFile, func(t *testing.T) {

			inputs, want := extractTestCase(t, testFile)
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
