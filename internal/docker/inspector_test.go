package docker

import (
	"context"
	"errors"
	"testing"

	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockClient is a mock implementation of the Docker client
type MockClient struct {
	mock.Mock
	client.Client
}

func (m *MockClient) ContainerList(ctx context.Context, options client.ContainerListOptions) (client.ContainerListResult, error) {
	args := m.Called(ctx, options)
	return args.Get(0).(client.ContainerListResult), args.Error(1)
}

func (m *MockClient) ContainerInspect(ctx context.Context, containerID string, options client.ContainerInspectOptions) (client.ContainerInspectResult, error) {
	args := m.Called(ctx, containerID, options)
	return args.Get(0).(client.ContainerInspectResult), args.Error(1)
}

func TestGetManagedContainers_Success(t *testing.T) {
	mockClient := new(MockClient)
	inspector := NewInspector(mockClient)

	// Mock data
	servicesDir := "/services"
	containerList := client.ContainerListResult{
		Items: []container.Summary{
			{
				ID:     "container1",
				Names:  []string{"/container1"},
				Image:  "image1",
				State:  "running",
				Status: "Up 5 minutes",
			},
			{
				ID:     "container2",
				Names:  []string{"/container2"},
				Image:  "image2",
				State:  "exited",
				Status: "Exited (0) 10 minutes ago",
			},
		},
	}

	containerInspect1 := client.ContainerInspectResult{
		Container: container.InspectResponse{
			Config: &container.Config{
				Labels: map[string]string{
					"com.docker.compose.project.working_dir": "/services/service1",
				},
			},
		},
	}

	containerInspect2 := client.ContainerInspectResult{
		Container: container.InspectResponse{
			Config: &container.Config{
				Labels: map[string]string{
					"com.docker.compose.project.working_dir": "/services/service2",
				},
			},
		},
	}

	// Setup expectations
	mockClient.On("ContainerList", mock.Anything, mock.Anything).Return(containerList, nil)
	mockClient.On("ContainerInspect", mock.Anything, "container1", mock.Anything).Return(containerInspect1, nil)
	mockClient.On("ContainerInspect", mock.Anything, "container2", mock.Anything).Return(containerInspect2, nil)

	// Call the method
	result, err := inspector.GetManagedContainers(servicesDir)

	// Assertions
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Contains(t, result, "service1")
	assert.Contains(t, result, "service2")
	assert.Len(t, result["service1"], 1)
	assert.Len(t, result["service2"], 1)
	assert.Equal(t, "container1", result["service1"][0].ID)
	assert.Equal(t, "container2", result["service2"][0].ID)

	// Verify that all expectations were met
	mockClient.AssertExpectations(t)
}

func TestGetManagedContainers_ContainerListError(t *testing.T) {
	mockClient := new(MockClient)
	inspector := NewInspector(mockClient)

	// Mock data
	servicesDir := "/services"
	expectedError := errors.New("failed to list containers")

	// Setup expectations
	mockClient.On("ContainerList", mock.Anything, mock.Anything).Return(client.ContainerListResult{}, expectedError)

	// Call the method
	result, err := inspector.GetManagedContainers(servicesDir)

	// Assertions
	assert.Error(t, err)
	assert.ErrorContains(t, err, expectedError.Error())
	assert.Nil(t, result)

	// Verify that all expectations were met
	mockClient.AssertExpectations(t)
}

func TestGetManagedContainers_ContainerInspectError(t *testing.T) {
	mockClient := new(MockClient)
	inspector := NewInspector(mockClient)

	// Mock data
	servicesDir := "/services"
	containerList := client.ContainerListResult{
		Items: []container.Summary{
			{
				ID:     "container1",
				Names:  []string{"/container1"},
				Image:  "image1",
				State:  "running",
				Status: "Up 5 minutes",
			},
		},
	}

	// Setup expectations
	mockClient.On("ContainerList", mock.Anything, mock.Anything).Return(containerList, nil)
	mockClient.On("ContainerInspect", mock.Anything, "container1", mock.Anything).Return(client.ContainerInspectResult{}, errors.New("failed to inspect container"))

	// Call the method
	result, err := inspector.GetManagedContainers(servicesDir)

	// Assertions
	assert.NoError(t, err)
	assert.Len(t, result, 0)

	// Verify that all expectations were met
	mockClient.AssertExpectations(t)
}

func TestGetManagedContainers_NoMatchingLabels(t *testing.T) {
	mockClient := MockClient{}
	inspector := NewInspector(&mockClient)

	// Mock data
	servicesDir := "/services"
	containerList := client.ContainerListResult{
		Items: []container.Summary{
			{
				ID:     "container1",
				Names:  []string{"/container1"},
				Image:  "image1",
				State:  "running",
				Status: "Up 5 minutes",
			},
		},
	}

	containerInspect := client.ContainerInspectResult{
		Container: container.InspectResponse{
			Config: &container.Config{
				Labels: map[string]string{
					"other.label": "value",
				},
			},
		},
	}

	// Setup expectations
	mockClient.On("ContainerList", mock.Anything, mock.Anything).Return(containerList, nil)
	mockClient.On("ContainerInspect", mock.Anything, "container1", mock.Anything).Return(containerInspect, nil)

	// Call the method
	result, err := inspector.GetManagedContainers(servicesDir)

	// Assertions
	assert.NoError(t, err)
	assert.Len(t, result, 0)

	// Verify that all expectations were met
	mockClient.AssertExpectations(t)
}

func TestGetManagedContainers_EmptyContainerList(t *testing.T) {
	mockClient := new(MockClient)
	inspector := NewInspector(mockClient)

	// Mock data
	servicesDir := "/services"
	containerList := client.ContainerListResult{
		Items: []container.Summary{},
	}

	// Setup expectations
	mockClient.On("ContainerList", mock.Anything, mock.Anything).Return(containerList, nil)

	// Call the method
	result, err := inspector.GetManagedContainers(servicesDir)

	// Assertions
	assert.NoError(t, err)
	assert.Len(t, result, 0)

	// Verify that all expectations were met
	mockClient.AssertExpectations(t)
}
