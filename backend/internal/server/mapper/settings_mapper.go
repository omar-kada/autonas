package mapper

import (
	"omar-kada/autonas/api"
	"omar-kada/autonas/models"
)

// SettingsMapper maps models.Config to api.Config
type SettingsMapper struct{}

// Map converts a models.Config to an api.Config
func (SettingsMapper) Map(settings models.Settings) api.Settings {
	return api.Settings{
		Branch: &settings.Branch,
		Cron:   &settings.CronPeriod,
		Repo:   settings.Repo,
	}
}

// UnMap transforms back from api.Config to models.Config
func (SettingsMapper) UnMap(settings api.Settings) models.Settings {
	res := models.Settings{
		Repo: settings.Repo,
	}
	if settings.Branch != nil {
		res.Branch = *settings.Branch
	}
	if settings.Cron != nil {
		res.CronPeriod = *settings.Cron
	}
	return res
}
