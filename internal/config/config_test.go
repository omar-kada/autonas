package config

import (
	"reflect"
	"testing"

	"github.com/elliotchance/orderedmap/v3"
)

func TestConfigPerService_BuildsCorrectArray(t *testing.T) {
	cfg := Config{
		Extra: map[string]any{
			"GLOBAL": "g",
		},
		Services: map[string]ServiceConfig{
			"svc": {
				Extra: map[string]any{
					"SVC_EXTRA": "s",
				},
			},
		},
	}

	got := cfg.PerService("svc")
	want := orderedmap.NewOrderedMapWithElements(
		&orderedmap.Element[string, string]{Key: "GLOBAL", Value: "g"},
		&orderedmap.Element[string, string]{Key: "SVC_EXTRA", Value: "s"},
	)
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("ConfigPerService mismatch\nwant=%#v\ngot =%#v", want, got)
	}
}
