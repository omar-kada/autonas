package server

import (
	"context"
	"omar-kada/autonas/api"
	"omar-kada/autonas/internal/process"
	"omar-kada/autonas/internal/storage"
	"time"
)

// Handler implements the generated strict server interface
type Handler struct {
	store   storage.DeploymentStorage
	manager process.Service
}

// NewHandler creates a new Handler
func NewHandler(store storage.DeploymentStorage, manager process.Service) *Handler {
	return &Handler{
		store:   store,
		manager: manager,
	}
}

// DeployementAPIList implements the StrictServerInterface interface
func (h *Handler) DeployementAPIList(_ context.Context, _ api.DeployementAPIListRequestObject) (api.DeployementAPIListResponseObject, error) {
	deps, err := h.store.GetDeployments()
	return api.DeployementAPIList200JSONResponse(deps), err
}

// DeployementAPIRead implements the StrictServerInterface interface
func (*Handler) DeployementAPIRead(_ context.Context, request api.DeployementAPIReadRequestObject) (api.DeployementAPIReadResponseObject, error) {
	// TODO: Implement your logic here
	// For now, we'll return a simple response
	return api.DeployementAPIRead200JSONResponse{
		Id:     request.Id,
		Title:  "Sample deployment",
		Time:   time.Now(),
		Diff:   "Sample diff",
		Status: "success",
	}, nil
}

// StatusAPIGet implements the StrictServerInterface interface
func (h *Handler) StatusAPIGet(_ context.Context, _ api.StatusAPIGetRequestObject) (api.StatusAPIGetResponseObject, error) {
	// TODO: Implement your logic here
	// For now, we'll return a simple response
	stacks, err := h.manager.GetManagerStacks()

	if err != nil {
		return nil, err
	}

	var result = make(map[string][]api.ContainerStatus)
	for stackName, containers := range stacks {
		for _, container := range containers {
			result[stackName] = append(result[stackName], api.ContainerStatus{
				ContainerId: container.ID,
				State:       api.ContainerStatusState(container.State),
				Name:        container.Name,
				Health:      api.ContainerStatusHealth(container.Health),
				StartedAt:   container.StartedAt,
			})
		}
	}
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

/*
// CreateHTTPHandler creates an HTTP handler for the API
func (h *Handler) CreateHTTPHandler() http.Handler {
	// Create the strict handler
	strictHandler := api.NewStrictHandler(h, nil)

	// Create the HTTP handler
	handler := api.Handler(strictHandler)

	// Add any middleware here if needed
	// For example:
	// handler = middleware1(handler)
	// handler = middleware2(handler)

	return handler
}
*/
