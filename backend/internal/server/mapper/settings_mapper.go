package mapper

import (
	"omar-kada/autonas/api"
	"omar-kada/autonas/models"
)

// SettingsMapper maps models.Settings to api.Settings
type SettingsMapper struct{}

// Map converts a models.Settings to an api.Settings
func (SettingsMapper) Map(settings models.Settings) api.Settings {
	token := settings.GetObfuscateToken()
	return api.Settings{
		Repo:     settings.Repo,
		Branch:   &settings.Branch,
		Cron:     &settings.Cron,
		Token:    &token,
		Username: &settings.Username,
	}
}

// UnMap transforms back from api.Settings to models.Settings
func (SettingsMapper) UnMap(settings api.Settings) models.Settings {
	res := models.Settings{
		Repo: settings.Repo,
	}
	if settings.Branch != nil {
		res.Branch = *settings.Branch
	}
	if settings.Cron != nil {
		res.Cron = *settings.Cron
	}
	if settings.Token != nil {
		res.Token = *settings.Token
	}
	if settings.Username != nil {
		res.Username = *settings.Username
	}
	return res
}
