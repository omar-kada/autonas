package config

import (
	"embed"
	"os"
	"path/filepath"
	"reflect"
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

	got := ConfigPerService(cfg, "svc")
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

func TestLoadConfig_SuccessWithOverride(t *testing.T) {
	f1 := getTempTestFile(t, "config.default.yaml")
	f2 := getTempTestFile(t, "config.override.yaml")

	cfg, err := LoadConfig([]string{f1, f2})
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}
	// Ensure merged values
	wantSvc := ServiceConfig{
		Port:    2000,
		Version: "v2",
		Extra:   map[string]any{"NEW_FIELD": "new_value"},
	}
	svc, ok := cfg.Services["svc"]
	if !ok {
		t.Fatalf("expected svc entry in services")
	}
	if !reflect.DeepEqual(svc, wantSvc) {
		t.Fatalf("unexpected svc values: want=%+v got=%+v", wantSvc, svc)
	}

}

func TestGetCurrentConfig_AfterLoadConfig(t *testing.T) {

	cfg, err := LoadConfig([]string{getTempTestFile(t, "config.default.yaml")})
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	// GetCurrentConfig should reflect last loaded cfg
	got := GetCurrentConfig()
	if !reflect.DeepEqual(got, cfg) {
		t.Fatalf("GetCurrentConfig mismatch\nwant=%#v\ngot =%#v", cfg, got)
	}
}

func TestLoadConfig_FileError(t *testing.T) {
	t.Run("missing file", func(t *testing.T) {
		if _, err := LoadConfig([]string{"/does/not/exist.yaml"}); err == nil {
			t.Fatalf("expected error for missing file")
		}
	})

	t.Run("invalid yaml", func(t *testing.T) {
		f := getTempTestFile(t, "config.invalid.yaml")
		if _, err := LoadConfig([]string{f}); err == nil {
			t.Fatalf("expected error for invalid yaml")
		}
	})
}
