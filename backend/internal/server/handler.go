package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"omar-kada/autonas/api"
	"omar-kada/autonas/internal/process"
	"omar-kada/autonas/internal/server/mapper"
	"omar-kada/autonas/internal/storage"
	"omar-kada/autonas/models"
)

// Handler implements the generated strict server interface
type Handler struct {
	store            storage.DeploymentStorage
	configStore      storage.ConfigStore
	processSvc       process.Service
	depMapper        mapper.DeploymentMapper
	depDetailsMapper mapper.DeploymentDetailsMapper
	diffMapper       mapper.DiffMapper
	statusMapper     mapper.StatusMapper
	statsMapper      mapper.StatsMapper
	configMapper     mapper.ConfigMapper
	settingsMapper   mapper.SettingsMapper
	featuresMapper   mapper.FeaturesMapper

	features models.Features
}

// NewHandler creates a new Handler
func NewHandler(store storage.DeploymentStorage, configStore storage.ConfigStore, service process.Service) *Handler {
	diffMapper := mapper.DiffMapper{}
	eventMapper := mapper.EventMapper{}

	return &Handler{
		store:            store,
		configStore:      configStore,
		processSvc:       service,
		depMapper:        mapper.NewDeploymentMapper(),
		depDetailsMapper: mapper.NewDeploymentDetailsMapper(diffMapper, eventMapper),
		diffMapper:       diffMapper,
		statusMapper:     mapper.StatusMapper{},
		statsMapper:      mapper.StatsMapper{},
		configMapper:     mapper.ConfigMapper{},
		features:         models.LoadFeatures(),
	}
}

// DeployementAPIList implements the StrictServerInterface interface
func (h *Handler) DeployementAPIList(_ context.Context, request api.DeployementAPIListRequestObject) (api.DeployementAPIListResponseObject, error) {
	offset, err := validateCursorOffset(request.Params.Offset)
	if err != nil {
		return nil, fmt.Errorf("invalid after value")
	}

	if request.Params.Limit <= 0 {
		return nil, fmt.Errorf("invalid first value")
	}

	deps, err := h.processSvc.GetDeployments(int(request.Params.Limit), offset)

	return api.DeployementAPIList200JSONResponse{
		Items:    models.ListMapper(h.depMapper.Map)(deps),
		PageInfo: h.depMapper.MapToPageInfo(deps, int(request.Params.Limit)),
	}, err
}

func validateCursorOffset(offsetStr *string) (uint64, error) {
	offset := uint64(0)
	var err error
	if offsetStr != nil && *offsetStr != "" {
		offset, err = strconv.ParseUint(*offsetStr, 10, 64)
	}
	return offset, err
}

// DeployementAPIRead implements the StrictServerInterface interface
func (h *Handler) DeployementAPIRead(_ context.Context, request api.DeployementAPIReadRequestObject) (api.DeployementAPIReadResponseObject, error) {
	id, err := strconv.ParseUint(request.Id, 10, 64)
	if err != nil {
		return nil, err
	}
	dep, err := h.store.GetDeployment(id)
	if err != nil {
		return nil, err
	} else if dep.ID == 0 {
		return api.DeployementAPIReaddefaultJSONResponse{
			Body: api.Error{
				Code:    http.StatusNotFound,
				Message: err.Error(),
			},
			StatusCode: http.StatusNotFound,
		}, err
	}

	return api.DeployementAPIRead200JSONResponse(h.depDetailsMapper.Map(dep)), err
}

// DeployementAPISync implements the StrictServerInterface interface
func (h *Handler) DeployementAPISync(_ context.Context, _ api.DeployementAPISyncRequestObject) (api.DeployementAPISyncResponseObject, error) {
	dep, err := h.processSvc.SyncDeployment()
	if err != nil {
		slog.Error(err.Error())
	}
	return api.DeployementAPISync200JSONResponse(h.depDetailsMapper.Map(dep)), err
}

// StatusAPIGet implements the StrictServerInterface interface
func (h *Handler) StatusAPIGet(_ context.Context, _ api.StatusAPIGetRequestObject) (api.StatusAPIGetResponseObject, error) {
	stacks, err := h.processSvc.GetManagedStacks()
	if err != nil {
		return nil, err
	}

	result := models.MapMapper[string](
		models.ListMapper(h.statusMapper.Map),
	)(stacks)

	var response []api.StackStatus
	for stackName, containers := range result {
		response = append(response, api.StackStatus{
			StackId:  stackName,
			Name:     stackName,
			Services: containers,
		})
	}
	return api.StatusAPIGet200JSONResponse(response), nil
}

// StatsAPIGet implements the StrictServerInterface interface
func (h *Handler) StatsAPIGet(_ context.Context, req api.StatsAPIGetRequestObject) (api.StatsAPIGetResponseObject, error) {
	stats, err := h.processSvc.GetCurrentStats(int(req.Days))
	if err != nil {
		return nil, err
	}
	return api.StatsAPIGet200JSONResponse(h.statsMapper.Map(stats)), nil
}

// DiffAPIGet implements the StrictServerInterface interface
func (h *Handler) DiffAPIGet(_ context.Context, _ api.DiffAPIGetRequestObject) (api.DiffAPIGetResponseObject, error) {
	fileDiffs, err := h.processSvc.GetDiff()
	if err != nil {
		return nil, err
	}
	return api.DiffAPIGet200JSONResponse(models.ListMapper(h.diffMapper.Map)(fileDiffs)), nil
}

// ConfigAPIGet implements the StrictServerInterface interface
func (h *Handler) ConfigAPIGet(_ context.Context, _ api.ConfigAPIGetRequestObject) (api.ConfigAPIGetResponseObject, error) {
	if !h.features.DisplayConfig {
		return api.ConfigAPIGetdefaultJSONResponse{
			Body: api.Error{
				Code:    http.StatusMethodNotAllowed,
				Message: "DISABLED",
			},
			StatusCode: http.StatusMethodNotAllowed,
		}, nil
	}
	config, err := h.configStore.Get()
	if err != nil {
		return nil, err
	}
	return api.ConfigAPIGet200JSONResponse(h.configMapper.Map(config)), nil
}

// ConfigAPISet implements the StrictServerInterface interface
func (h *Handler) ConfigAPISet(_ context.Context, r api.ConfigAPISetRequestObject) (api.ConfigAPISetResponseObject, error) {
	if !h.features.EditConfig {
		return api.ConfigAPISetdefaultJSONResponse{
			Body: api.Error{
				Code:    http.StatusMethodNotAllowed,
				Message: "DISABLED",
			},
			StatusCode: http.StatusMethodNotAllowed,
		}, nil
	}
	config := h.configMapper.UnMap(api.Config(*r.Body))
	oldConfig, err := h.configStore.Get()
	if err != nil {
		return nil, err
	}
	oldConfig.Environment = config.Environment
	oldConfig.Services = config.Services
	err = h.configStore.Update(oldConfig)
	if err != nil {
		return nil, err
	}
	return api.ConfigAPISet200JSONResponse(h.configMapper.Map(oldConfig)), nil
}

// SettingsAPIGet implements the StrictServerInterface interface
func (h *Handler) SettingsAPIGet(_ context.Context, _ api.SettingsAPIGetRequestObject) (api.SettingsAPIGetResponseObject, error) {
	config, err := h.configStore.Get()
	if err != nil {
		return nil, err
	}
	return api.SettingsAPIGet200JSONResponse(h.settingsMapper.Map(config.Settings)), nil
}

// SettingsAPISet implements the StrictServerInterface interface
func (h *Handler) SettingsAPISet(_ context.Context, r api.SettingsAPISetRequestObject) (api.SettingsAPISetResponseObject, error) {
	if !h.features.EditSettings {
		return api.SettingsAPISetdefaultJSONResponse{
			Body: api.Error{
				Code:    http.StatusMethodNotAllowed,
				Message: "DISABLED",
			},
			StatusCode: http.StatusMethodNotAllowed,
		}, nil
	}
	oldConfig, err := h.configStore.Get()
	if err != nil {
		return nil, err
	}
	settings := h.settingsMapper.UnMap(api.Settings(*r.Body))

	if models.IsObfuscated(settings.Token) {
		settings.Token = oldConfig.Settings.Token // keep old token when obfuscated
	}
	oldConfig.Settings = settings
	err = h.configStore.Update(oldConfig)
	if err != nil {
		return nil, err
	}
	return api.SettingsAPISet200JSONResponse(h.settingsMapper.Map(settings)), nil
}

// FeaturesAPIGet implements the StrictServerInterface interface
func (h *Handler) FeaturesAPIGet(_ context.Context, _ api.FeaturesAPIGetRequestObject) (api.FeaturesAPIGetResponseObject, error) {
	return api.FeaturesAPIGet200JSONResponse(h.featuresMapper.Map(h.features)), nil
}
