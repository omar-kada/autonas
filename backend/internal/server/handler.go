package server

import (
	"context"
	"fmt"
	"log/slog"
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
	processSvc       process.Service
	depMapper        mapper.DeploymentMapper
	depDetailsMapper mapper.DeploymentDetailsMapper
	diffMapper       mapper.DiffMapper
	statusMapper     mapper.StatusMapper
	statsMapper      mapper.StatsMapper
}

// NewHandler creates a new Handler
func NewHandler(store storage.DeploymentStorage, service process.Service) *Handler {
	diffMapper := mapper.DiffMapper{}
	eventMapper := mapper.EventMapper{}

	return &Handler{
		store:            store,
		processSvc:       service,
		depMapper:        mapper.NewDeploymentMapper(),
		depDetailsMapper: mapper.NewDeploymentDetailsMapper(diffMapper, eventMapper),
		diffMapper:       diffMapper,
		statusMapper:     mapper.StatusMapper{},
		statsMapper:      mapper.StatsMapper{},
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
				Code:    404,
				Message: err.Error(),
			},
			StatusCode: 404,
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
