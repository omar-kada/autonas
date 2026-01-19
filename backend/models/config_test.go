package models

import (
	"reflect"
	"testing"

	"github.com/elliotchance/orderedmap/v3"
	"github.com/stretchr/testify/assert"
)

func TestConfigPerService_BuildsCorrectArray(t *testing.T) {
	cfg := Config{
		Environment: map[string]string{
			"GLOBAL": "g",
		},
		Services: map[string]map[string]string{
			"svc": {
				"SVC_EXTRA": "s",
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

func TestGetEnabledServices_FiltersCorrectly(t *testing.T) {
	cfg := Config{
		Environment: map[string]string{
			"GLOBAL": "g",
		},
		Services: map[string]map[string]string{
			"svc": {
				"SVC_EXTRA": "s",
			},
			"svc2": {
				"SVC_EXTRA": "s",
			},
		},
	}

	want := []string{"svc"}
	assert.EqualValues(t, want, cfg.GetEnabledServices())
}
