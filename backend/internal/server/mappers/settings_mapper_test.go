package mappers

import (
	"testing"

	"omar-kada/autonas/api"
	"omar-kada/autonas/models"

	"github.com/stretchr/testify/assert"
)

func TestSettingsMapper_Map(t *testing.T) {
	main := "main"
	cron := "0 0 * * *"
	username := "user"
	token := "123456789123456789"
	obfuscatedToken := models.ObfuscateToken(token)
	empty := ""
	cases := []struct {
		name string
		in   models.Settings
		want api.Settings
	}{
		{
			name: "basic",
			in: models.Settings{
				Repo:     "https://github.com/example/repo",
				Branch:   main,
				Cron:     cron,
				Username: username,
				Token:    token,
			},
			want: api.Settings{
				Repo:     "https://github.com/example/repo",
				Branch:   &main,
				Cron:     &cron,
				Token:    &obfuscatedToken,
				Username: &username,
			},
		},
		{
			name: "empty",
			in: models.Settings{
				Repo: "",
			},
			want: api.Settings{
				Repo:     "",
				Branch:   &empty,
				Cron:     &empty,
				Token:    &empty,
				Username: &empty,
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
				Repo:   repo,
				Branch: branch,
				Cron:   cron,
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
				Repo:   "",
				Branch: "",
				Cron:   "",
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
