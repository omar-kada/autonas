// Package mappers provides functionality for mapping between different data models.

package mappers

import (
	"testing"

	"omar-kada/autonas/api"
	"omar-kada/autonas/models"

	"github.com/stretchr/testify/assert"
)

func TestConfigMapper_Map(t *testing.T) {
	cases := []struct {
		name string
		in   models.Config
		want api.Config
	}{
		{
			name: "basic",
			in: models.Config{
				Settings: models.Settings{
					Repo:   "https://github.com/example/repo",
					Branch: "main",
					Cron:   "0 0 * * *",
				},
				Environment: models.Environment{
					"key1": "value1",
					"key2": "value2",
				},
				Services: map[string]models.ServiceConfig{
					"service1": {
						"key1": "value1",
						"key2": "value2",
					},
					"service2": {
						"key3": "value3",
						"key4": "value4",
					},
				},
			},
			want: api.Config{
				GlobalVariables: map[string]string{
					"key1": "value1",
					"key2": "value2",
				},
				Services: map[string]map[string]string{
					"service1": {
						"key1": "value1",
						"key2": "value2",
					},
					"service2": {
						"key3": "value3",
						"key4": "value4",
					},
				},
			},
		},
		{
			name: "empty",
			in: models.Config{
				Settings:    models.Settings{},
				Environment: models.Environment{},
				Services:    map[string]models.ServiceConfig{},
			},
			want: api.Config{
				GlobalVariables: map[string]string{},
				Services:        map[string]map[string]string{},
			},
		},
	}

	m := ConfigMapper{}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := m.Map(tc.in)
			assert.Equal(t, tc.want, got)
		})
	}
}
