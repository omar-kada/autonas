package mappers

import (
	"testing"
	"time"

	"omar-kada/autonas/api"
	"omar-kada/autonas/models"

	"github.com/moby/moby/api/types/container"
	"github.com/stretchr/testify/assert"
)

func TestStatsMapper_Map(t *testing.T) {
	now := time.Now().Truncate(time.Second)
	next := now.Add(24 * time.Hour)

	cases := []struct {
		name string
		in   models.Stats
		want api.Stats
	}{
		{
			name: "basic",
			in: models.Stats{
				Author:     "alice",
				Error:      1,
				Success:    2,
				LastDeploy: now,
				NextDeploy: next,
				LastStatus: models.DeploymentStatusRunning,
				Health:     container.HealthStatus("healthy"),
			},
			want: api.Stats{
				Author:     "alice",
				Error:      1,
				Success:    2,
				LastDeploy: now,
				NextDeploy: next,
				Status:     api.DeploymentStatus(models.DeploymentStatusRunning),
				Health:     api.ContainerHealth("healthy"),
			},
		},
		{
			name: "zero-times-empty-health",
			in: models.Stats{
				Author:     "bob",
				Error:      0,
				Success:    0,
				LastDeploy: time.Time{},
				NextDeploy: time.Time{},
				LastStatus: models.DeploymentStatusPlanned,
				Health:     container.HealthStatus("none"),
			},
			want: api.Stats{
				Author:     "bob",
				Error:      0,
				Success:    0,
				LastDeploy: time.Time{},
				NextDeploy: time.Time{},
				Status:     api.DeploymentStatus(models.DeploymentStatusPlanned),
				Health:     api.ContainerHealth(container.HealthStatus("none")),
			},
		},
	}

	m := StatsMapper{}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := m.Map(tc.in)
			assert.Equal(t, tc.want, got)
		})
	}
}
