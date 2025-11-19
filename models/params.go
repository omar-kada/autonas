package models

// RunParams represent the global parameters of the run command
type RunParams struct {
	ConfigFiles  []string
	Repo         string
	Branch       string
	WorkingDir   string
	ServicesDir  string
	CronPeriod   string
	AddWritePerm bool
	Port         int
}
