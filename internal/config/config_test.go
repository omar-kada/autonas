package config

import (
	"reflect"
	"testing"
)

func TestConfigPerService_BuildsCorrectArray(t *testing.T) {
	cfg := Config{

		Extra: map[string]any{
			"GLOBAL": "g",
		},
		Services: map[string]ServiceConfig{
			"svc": {
				Port:    8080,
				Version: "v1",
				Extra:   map[string]any{"SVC_EXTRA": "s"},
			},
		},
	}

	got := cfg.PerService("svc")
	want := []Variable{
		{Key: "GLOBAL", Value: "g"},
		{Key: "PORT", Value: "8080"},
		{Key: "VERSION", Value: "v1"},
		{Key: "SVC_EXTRA", Value: "s"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("ConfigPerService mismatch\nwant=%#v\ngot =%#v", want, got)
	}
}
