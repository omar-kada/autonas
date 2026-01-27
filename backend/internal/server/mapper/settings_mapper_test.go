package mapper

import (
	"testing"

	"omar-kada/autonas/api"
	"omar-kada/autonas/models"

	"github.com/stretchr/testify/assert"
)

func TestSettingsMapper_Map(t *testing.T) {
	main := "main"
	cron := "0 0 * * *"
	empty := ""
	cases := []struct {
		name string
		in   models.Settings
		want api.Settings
	}{
		{
			name: "basic",
			in: models.Settings{
				Repo:       "https://github.com/example/repo",
				Branch:     "main",
				CronPeriod: "0 0 * * *",
			},
			want: api.Settings{
				Repo:   "https://github.com/example/repo",
				Branch: &main,
				Cron:   &cron,
			},
		},
		{
			name: "empty",
			in: models.Settings{
				Repo:       "",
				Branch:     empty,
				CronPeriod: empty,
			},
			want: api.Settings{
				Repo:   "",
				Branch: &empty,
				Cron:   &empty,
			},
		},
	}

	m := SettingsMapper{}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := m.Map(tc.in)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestSettingsMapper_UnMap(t *testing.T) {
	branch := "main"
	cron := "0 0 * * *"
	repo := "https://github.com/example/repo"

	cases := []struct {
		name string
		in   api.Settings
		want models.Settings
	}{
		{
			name: "basic",
			in: api.Settings{
				Repo:   repo,
				Branch: &branch,
				Cron:   &cron,
			},
			want: models.Settings{
				Repo:       repo,
				Branch:     branch,
				CronPeriod: cron,
			},
		},
		{
			name: "empty",
			in: api.Settings{
				Branch: nil,
				Cron:   nil,
				Repo:   "",
			},
			want: models.Settings{
				Repo:       "",
				Branch:     "",
				CronPeriod: "",
			},
		},
	}

	m := SettingsMapper{}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := m.UnMap(tc.in)
			assert.Equal(t, tc.want, got)
		})
	}
}
