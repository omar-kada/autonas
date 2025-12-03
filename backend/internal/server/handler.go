package server

import (
	"context"
	"omar-kada/autonas/api"
	"omar-kada/autonas/internal/storage"
	"time"
)

// Handler implements the generated strict server interface
type Handler struct {
	store storage.Storage
}

// NewHandler creates a new Handler
func NewHandler(store storage.Storage) *Handler {
	return &Handler{
		store: store,
	}
}

// DeployementAPIList implements the StrictServerInterface interface
func (h *Handler) DeployementAPIList(ctx context.Context, request api.DeployementAPIListRequestObject) (api.DeployementAPIListResponseObject, error) {
	deps, err := h.store.GetDeployments()
	return api.DeployementAPIList200JSONResponse(deps), err
}

// DeployementAPIRead implements the StrictServerInterface interface
func (h *Handler) DeployementAPIRead(ctx context.Context, request api.DeployementAPIReadRequestObject) (api.DeployementAPIReadResponseObject, error) {
	// TODO: Implement your logic here
	// For now, we'll return a simple response
	return api.DeployementAPIRead200JSONResponse{
		Id:     request.Id,
		Title:  "Sample deployment",
		Time:   time.Now(),
		Diff:   "Sample diff",
		Status: "success",
		Logs:   []string{"Deployment started", "Deployment completed"},
	}, nil
}

// StatusAPIGet implements the StrictServerInterface interface
func (h *Handler) StatusAPIGet(ctx context.Context, request api.StatusAPIGetRequestObject) (api.StatusAPIGetResponseObject, error) {
	// TODO: Implement your logic here
	// For now, we'll return a simple response
	return api.StatusAPIGet200JSONResponse{
		{
			StackId: "stack1",
			Name:    "MyStack",
			Services: []api.ContainerStatus{
				{
					ContainerId: "container1",
					State:       "running",
					Name:        "service1",
					Health:      "healthy",
					CreatedAt:   time.Now(),
				},
				{
					ContainerId: "container2",
					State:       "stopped",
					Name:        "service2",
					Health:      "unhealthy",
					CreatedAt:   time.Now(),
				},
			},
		},
	}, nil
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
