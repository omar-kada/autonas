package server

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"omar-kada/autonas/api"
	"omar-kada/autonas/internal/storage"
	"omar-kada/autonas/models"

	"github.com/moby/moby/api/types/container"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockProcess struct {
	mock.Mock
}

func (m *MockProcess) SyncDeployment() (models.Deployment, error) {
	args := m.Called()
	return args.Get(0).(models.Deployment), args.Error(1)
}

func (m *MockProcess) GetCurrentStats(days int) (models.Stats, error) {
	args := m.Called(days)
	return args.Get(0).(models.Stats), args.Error(1)
}

func (m *MockProcess) GetDiff() ([]models.FileDiff, error) {
	args := m.Called()
	return args.Get(0).([]models.FileDiff), args.Error(1)
}

func (m *MockProcess) GetManagedStacks() (map[string][]models.ContainerSummary, error) {
	args := m.Called()
	return args.Get(0).(map[string][]models.ContainerSummary), args.Error(1)
}

func (m *MockProcess) GetDeployments(limit int, offset uint64) ([]models.Deployment, error) {
	args := m.Called(limit, offset)
	return args.Get(0).([]models.Deployment), args.Error(1)
}

type MockStore struct {
	mock.Mock
}

func (m *MockStore) GetDeployments(c storage.Cursor[uint64]) ([]models.Deployment, error) {
	args := m.Called(c)
	return args.Get(0).([]models.Deployment), args.Error(1)
}

func (m *MockStore) GetDeployment(id uint64) (models.Deployment, error) {
	args := m.Called(id)
	return args.Get(0).(models.Deployment), args.Error(1)
}

func (m *MockStore) InitDeployment(title string, author string, diff string, files []models.FileDiff) (models.Deployment, error) {
	args := m.Called(title, author, diff, files)
	return args.Get(0).(models.Deployment), args.Error(1)
}

func (m *MockStore) EndDeployment(deploymentID uint64, status models.DeploymentStatus) error {
	args := m.Called(deploymentID, status)
	return args.Error(0)
}

func (m *MockStore) GetLastDeployment() (models.Deployment, error) {
	args := m.Called()
	return args.Get(0).(models.Deployment), args.Error(1)
}

func (m *MockStore) Get() (models.Config, error) {
	args := m.Called()
	return args.Get(0).(models.Config), args.Error(1)
}

func (m *MockStore) Update(config models.Config) error {
	args := m.Called(config)
	return args.Error(0)
}

func (m *MockStore) SetOnChange(fn func(models.Config, models.Config)) {
	m.Called(fn)
}

func TestDeployementAPIList_Success(t *testing.T) {
	m := &MockProcess{}
	store := &MockStore{}
	h := NewHandler(store, store, m)

	deps := []models.Deployment{
		{ID: 1, Title: "first", Author: "alice", Diff: "d1", Status: models.DeploymentStatusSuccess},
		{ID: 2, Title: "second", Author: "bob", Diff: "d2", Status: models.DeploymentStatusRunning},
	}
	m.On("GetDeployments", 2, uint64(0)).Return(deps, nil)

	req := api.DeployementAPIListRequestObject{Params: api.DeployementAPIListParams{Limit: 2}}
	resp, err := h.DeployementAPIList(context.Background(), req)
	assert.NoError(t, err)

	switch r := resp.(type) {
	case api.DeployementAPIList200JSONResponse:
		assert.Equal(t, 2, len(r.Items))
		assert.Equal(t, "2", r.PageInfo.EndCursor)
	default:
		t.Fatalf("unexpected response type: %T", resp)
	}

	m.AssertExpectations(t)
}

func TestDeployementAPIList_InvalidOffset(t *testing.T) {
	m := &MockProcess{}
	store := &MockStore{}
	h := NewHandler(store, store, m)

	off := "notuint"
	req := api.DeployementAPIListRequestObject{Params: api.DeployementAPIListParams{Limit: 1, Offset: &off}}
	resp, err := h.DeployementAPIList(context.Background(), req)
	assert.Nil(t, resp)
	assert.EqualError(t, err, "invalid after value")
}

func TestDeployementAPIRead_Success(t *testing.T) {
	m := &MockProcess{}
	store := &MockStore{}
	h := NewHandler(store, store, m)

	dep := models.Deployment{ID: 10, Title: "Automatic Deploy", Author: "ci", Diff: "diff"}
	store.On("GetDeployment", uint64(10)).Return(dep, nil)

	req := api.DeployementAPIReadRequestObject{Id: "10"}
	resp, err := h.DeployementAPIRead(context.Background(), req)
	assert.NoError(t, err)

	switch r := resp.(type) {
	case api.DeployementAPIRead200JSONResponse:
		assert.Equal(t, "Automatic Deploy", r.Title)
		assert.Equal(t, "diff", r.Diff)
	default:
		t.Fatalf("unexpected response type: %T", resp)
	}

	store.AssertExpectations(t)
}

func TestDeployementAPIRead_InvalidID(t *testing.T) {
	m := &MockProcess{}
	store := &MockStore{}
	h := NewHandler(store, store, m)

	req := api.DeployementAPIReadRequestObject{Id: "abc"}
	resp, err := h.DeployementAPIRead(context.Background(), req)
	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestDeployementAPISync_SuccessAndError(t *testing.T) {
	m := &MockProcess{}
	store := &MockStore{}
	h := NewHandler(store, store, m)

	dep := models.Deployment{ID: 99, Title: "Automatic Deploy", Author: "ci", Diff: "dd", Status: models.DeploymentStatusRunning}
	m.On("SyncDeployment").Return(dep, nil)

	resp, err := h.DeployementAPISync(context.Background(), api.DeployementAPISyncRequestObject{})
	assert.NoError(t, err)

	switch r := resp.(type) {
	case api.DeployementAPISync200JSONResponse:
		assert.Equal(t, "Automatic Deploy", r.Title)
		assert.Equal(t, "dd", r.Diff)
	default:
		t.Fatalf("unexpected resp type: %T", resp)
	}

	// now return an error (handler should return both response and error)
	errTest := errors.New("sync failed")
	m.ExpectedCalls = nil
	m.On("SyncDeployment").Return(models.Deployment{}, errTest)

	_, err2 := h.DeployementAPISync(context.Background(), api.DeployementAPISyncRequestObject{})
	assert.Equal(t, errTest, err2)

	m.AssertExpectations(t)
}

func TestStatusAPIGet_Success(t *testing.T) {
	m := &MockProcess{}
	store := &MockStore{}
	h := NewHandler(store, store, m)

	stacks := map[string][]models.ContainerSummary{
		"stack1": {{ID: "c1", Name: "c1", Image: "img1", State: container.StateRunning, Health: container.Healthy}},
	}
	m.On("GetManagedStacks").Return(stacks, nil)

	resp, err := h.StatusAPIGet(context.Background(), api.StatusAPIGetRequestObject{})
	assert.NoError(t, err)

	switch r := resp.(type) {
	case api.StatusAPIGet200JSONResponse:
		assert.Equal(t, 1, len(r))
		assert.Equal(t, "stack1", r[0].StackId)
	default:
		t.Fatalf("unexpected resp type: %T", resp)
	}

	m.AssertExpectations(t)
}

func TestStatsAPIGet_Success(t *testing.T) {
	m := &MockProcess{}
	store := &MockStore{}
	h := NewHandler(store, store, m)

	next := time.Now().Add(1 * time.Hour)
	stats := models.Stats{Author: "bob", Error: 1, Success: 2, LastStatus: models.DeploymentStatusError, NextDeploy: next}
	m.On("GetCurrentStats", 7).Return(stats, nil)

	req := api.StatsAPIGetRequestObject{Days: 7}
	resp, err := h.StatsAPIGet(context.Background(), req)
	assert.NoError(t, err)

	switch r := resp.(type) {
	case api.StatsAPIGet200JSONResponse:
		assert.Equal(t, int32(2), r.Success)
		assert.Equal(t, int32(1), r.Error)
		assert.Equal(t, "bob", r.Author)
		assert.Equal(t, next, r.NextDeploy)
	default:
		t.Fatalf("unexpected resp type: %T", resp)
	}

	m.AssertExpectations(t)
}

func TestDiffAPIGet_Success(t *testing.T) {
	m := &MockProcess{}
	store := &MockStore{}
	h := NewHandler(store, store, m)

	fileDiffs := []models.FileDiff{{OldFile: "file1.txt", NewFile: "file1.txt", Diff: "d1"}}
	m.On("GetDiff").Return(fileDiffs, nil)

	resp, err := h.DiffAPIGet(context.Background(), api.DiffAPIGetRequestObject{})
	assert.NoError(t, err)

	switch r := resp.(type) {
	case api.DiffAPIGet200JSONResponse:
		assert.Equal(t, 1, len(r))
		assert.Equal(t, "file1.txt", r[0].OldFile)
	default:
		t.Fatalf("unexpected resp type: %T", resp)
	}

	m.AssertExpectations(t)
}

func TestDeployementAPIList_GetDeploymentsError(t *testing.T) {
	m := &MockProcess{}
	store := &MockStore{}
	h := NewHandler(store, store, m)

	errList := errors.New("db error")
	m.On("GetDeployments", 2, uint64(0)).Return([]models.Deployment{}, errList)

	req := api.DeployementAPIListRequestObject{Params: api.DeployementAPIListParams{Limit: 2}}
	resp, err := h.DeployementAPIList(context.Background(), req)
	assert.Error(t, err)
	assert.Equal(t, errList, err)

	switch r := resp.(type) {
	case api.DeployementAPIList200JSONResponse:
		assert.Equal(t, 0, len(r.Items))
	default:
		t.Fatalf("unexpected response type: %T", resp)
	}

	m.AssertExpectations(t)
}

func TestDeployementAPIList_InvalidLimit(t *testing.T) {
	m := &MockProcess{}
	store := &MockStore{}
	h := NewHandler(store, store, m)

	req := api.DeployementAPIListRequestObject{Params: api.DeployementAPIListParams{Limit: 0}}
	resp, err := h.DeployementAPIList(context.Background(), req)
	assert.Nil(t, resp)
	assert.EqualError(t, err, "invalid first value")
}

func TestDeployementAPIRead_GetDeploymentError(t *testing.T) {
	m := &MockProcess{}
	store := &MockStore{}
	h := NewHandler(store, store, m)

	errGet := errors.New("not found")
	store.On("GetDeployment", uint64(10)).Return(models.Deployment{}, errGet)

	req := api.DeployementAPIReadRequestObject{Id: "10"}
	resp, err := h.DeployementAPIRead(context.Background(), req)
	assert.Error(t, err)
	assert.Equal(t, errGet, err)
	assert.Nil(t, resp)

	store.AssertExpectations(t)
}

func TestStatusAPIGet_Error(t *testing.T) {
	m := &MockProcess{}
	store := &MockStore{}
	h := NewHandler(store, store, m)

	errStacks := errors.New("failed to get stacks")
	m.On("GetManagedStacks").Return(map[string][]models.ContainerSummary{}, errStacks)

	resp, err := h.StatusAPIGet(context.Background(), api.StatusAPIGetRequestObject{})
	assert.Nil(t, resp)
	assert.Equal(t, errStacks, err)

	m.AssertExpectations(t)
}

func TestStatsAPIGet_Error(t *testing.T) {
	m := &MockProcess{}
	store := &MockStore{}
	h := NewHandler(store, store, m)

	errStats := errors.New("stats error")
	m.On("GetCurrentStats", 7).Return(models.Stats{}, errStats)

	req := api.StatsAPIGetRequestObject{Days: 7}
	resp, err := h.StatsAPIGet(context.Background(), req)
	assert.Nil(t, resp)
	assert.Equal(t, errStats, err)

	m.AssertExpectations(t)
}

func TestDiffAPIGet_Error(t *testing.T) {
	m := &MockProcess{}
	store := &MockStore{}
	h := NewHandler(store, store, m)

	errDiff := errors.New("diff error")
	m.On("GetDiff").Return([]models.FileDiff{}, errDiff)

	resp, err := h.DiffAPIGet(context.Background(), api.DiffAPIGetRequestObject{})
	assert.Nil(t, resp)
	assert.Equal(t, errDiff, err)

	m.AssertExpectations(t)
}

func TestConfigAPIGet_Success(t *testing.T) {
	m := &MockProcess{}
	store := &MockStore{}
	h := NewHandler(store, store, m)
	h.features.DisplayConfig = true

	config := models.Config{
		Environment: models.Environment{
			"ENV": "VALUE",
		},
	}
	store.On("Get").Return(config, nil)

	resp, err := h.ConfigAPIGet(context.Background(), api.ConfigAPIGetRequestObject{})
	assert.NoError(t, err)

	switch r := resp.(type) {
	case api.ConfigAPIGet200JSONResponse:
		assert.Equal(t, "VALUE", r.GlobalVariables["ENV"])
	default:
		t.Fatalf("unexpected resp type: %T", resp)
	}

	store.AssertExpectations(t)
}

func TestConfigAPIGet_Disabled(t *testing.T) {
	m := &MockProcess{}
	store := &MockStore{}
	h := NewHandler(store, store, m)
	h.features.DisplayConfig = false

	resp, err := h.ConfigAPIGet(context.Background(), api.ConfigAPIGetRequestObject{})
	assert.NoError(t, err)

	switch r := resp.(type) {
	case api.ConfigAPIGetdefaultJSONResponse:
		assert.Equal(t, http.StatusMethodNotAllowed, r.StatusCode)
		assert.Equal(t, "DISABLED", r.Body.Message)
	default:
		t.Fatalf("unexpected resp type: %T", resp)
	}
}

func TestConfigAPIGet_Error(t *testing.T) {
	m := &MockProcess{}
	store := &MockStore{}
	h := NewHandler(store, store, m)
	h.features.DisplayConfig = true

	errConfig := errors.New("config error")
	store.On("Get").Return(models.Config{}, errConfig)

	resp, err := h.ConfigAPIGet(context.Background(), api.ConfigAPIGetRequestObject{})
	assert.Nil(t, resp)
	assert.Equal(t, errConfig, err)

	store.AssertExpectations(t)
}

func TestFeaturesAPIGet_Success(t *testing.T) {
	m := &MockProcess{}
	store := &MockStore{}
	h := NewHandler(store, store, m)

	resp, err := h.FeaturesAPIGet(context.Background(), api.FeaturesAPIGetRequestObject{})
	assert.NoError(t, err)

	switch r := resp.(type) {
	case api.FeaturesAPIGet200JSONResponse:
		assert.Equal(t, h.features.DisplayConfig, r.DisplayConfig)
	default:
		t.Fatalf("unexpected resp type: %T", resp)
	}
}
func TestSettingsAPIGet_Success(t *testing.T) {
	m := &MockProcess{}
	store := &MockStore{}
	h := NewHandler(store, store, m)

	settings := models.Settings{
		Repo:       "test-repo",
		Branch:     "main",
		CronPeriod: "0 0 * * *",
	}
	store.On("Get").Return(models.Config{Settings: settings}, nil)

	resp, err := h.SettingsAPIGet(context.Background(), api.SettingsAPIGetRequestObject{})
	assert.NoError(t, err)

	switch r := resp.(type) {
	case api.SettingsAPIGet200JSONResponse:
		assert.Equal(t, "test-repo", r.Repo)
		assert.Equal(t, "main", *r.Branch)
		assert.Equal(t, "0 0 * * *", *r.Cron)
	default:
		t.Fatalf("unexpected resp type: %T", resp)
	}

	store.AssertExpectations(t)
}

func TestSettingsAPIGet_Error(t *testing.T) {
	m := &MockProcess{}
	store := &MockStore{}
	h := NewHandler(store, store, m)

	errSettings := errors.New("settings error")
	store.On("Get").Return(models.Config{}, errSettings)

	resp, err := h.SettingsAPIGet(context.Background(), api.SettingsAPIGetRequestObject{})
	assert.Nil(t, resp)
	assert.Equal(t, errSettings, err)

	store.AssertExpectations(t)
}

func TestSettingsAPISet_Success(t *testing.T) {
	m := &MockProcess{}
	store := &MockStore{}
	h := NewHandler(store, store, m)
	h.features.EditSettings = true

	oldConfig := models.Config{
		Environment: models.Environment{
			"ENV": "VALUE",
		},
		Services: map[string]models.ServiceConfig{
			"service1": {"key1": "value1"},
		},
		Settings: models.Settings{
			Repo:       "old-repo",
			Branch:     "old-branch",
			CronPeriod: "old-cron",
		},
	}
	newSettings := api.Settings{
		Repo:   "new-repo",
		Branch: ptr("new-branch"),
		Cron:   ptr("new-cron"),
	}

	store.On("Get").Return(oldConfig, nil)
	store.On("Update", mock.MatchedBy(func(newCfg models.Config) bool {
		// Check that only settings are updated
		assert.Equal(t, oldConfig.Environment, newCfg.Environment)
		assert.Equal(t, oldConfig.Services, newCfg.Services)
		assert.Equal(t, models.Settings{
			Repo:       newSettings.Repo,
			Branch:     *newSettings.Branch,
			CronPeriod: *newSettings.Cron,
		}, newCfg.Settings)
		return true
	})).Return(nil)

	req := api.SettingsAPISetRequestObject{Body: &newSettings}
	resp, err := h.SettingsAPISet(context.Background(), req)
	assert.NoError(t, err)

	switch r := resp.(type) {
	case api.SettingsAPISet200JSONResponse:
		assert.Equal(t, "new-repo", r.Repo)
		assert.Equal(t, "new-branch", *r.Branch)
		assert.Equal(t, "new-cron", *r.Cron)
	default:
		t.Fatalf("unexpected resp type: %T", resp)
	}

	store.AssertExpectations(t)
}

func TestSettingsAPISet_Disabled(t *testing.T) {
	m := &MockProcess{}
	store := &MockStore{}
	h := NewHandler(store, store, m)
	h.features.EditSettings = false

	settings := api.Settings{
		Repo:   "test-repo",
		Branch: ptr("main"),
		Cron:   ptr("0 0 * * *"),
	}

	req := api.SettingsAPISetRequestObject{Body: &settings}
	resp, err := h.SettingsAPISet(context.Background(), req)
	assert.NoError(t, err)

	switch r := resp.(type) {
	case api.SettingsAPISetdefaultJSONResponse:
		assert.Equal(t, http.StatusMethodNotAllowed, r.StatusCode)
		assert.Equal(t, "DISABLED", r.Body.Message)
	default:
		t.Fatalf("unexpected resp type: %T", resp)
	}
}

func TestSettingsAPISet_Error(t *testing.T) {
	m := &MockProcess{}
	store := &MockStore{}
	h := NewHandler(store, store, m)
	h.features.EditSettings = true

	settings := api.Settings{
		Repo:   "test-repo",
		Branch: ptr("main"),
		Cron:   ptr("0 0 * * *"),
	}

	errSettings := errors.New("settings error")
	store.On("Get").Return(models.Config{}, errSettings)

	req := api.SettingsAPISetRequestObject{Body: &settings}
	resp, err := h.SettingsAPISet(context.Background(), req)
	assert.Nil(t, resp)
	assert.Equal(t, errSettings, err)

	store.AssertExpectations(t)
}

func ptr[T any](v T) *T {
	return &v
}
