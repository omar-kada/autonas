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

// MockClient is a mock implementation of the Client interface
type MockClient struct {
	mock.Mock
}

func (m *MockClient) ContainerList(ctx context.Context, options client.ContainerListOptions) (client.ContainerListResult, error) {
	args := m.Called(ctx, options)
	return args.Get(0).(client.ContainerListResult), args.Error(1)
}

func (m *MockClient) ContainerInspect(ctx context.Context, containerID string, options client.ContainerInspectOptions) (client.ContainerInspectResult, error) {
	args := m.Called(ctx, containerID, options)
	return args.Get(0).(client.ContainerInspectResult), args.Error(1)
}

func TestGetManagedContainers(t *testing.T) {
	mockClient := new(MockClient)

	// Test successful case
	mockClient.On("ContainerList", mock.Anything, mock.Anything).Once().Return(client.ContainerListResult{
		Items: []container.Summary{
			{
				ID:     "container1",
				Names:  []string{"/container1"},
				Image:  "image1",
				State:  "running",
				Status: "Up 1 hour",
			},
			{
				ID:     "container2",
				Names:  []string{"/container2"},
				Image:  "image2",
				State:  "exited",
				Status: "Exited (0) 2 hours ago",
			},
		},
	}, nil)

	mockClient.On("ContainerInspect", mock.Anything, "container1", mock.Anything).Return(client.ContainerInspectResult{
		Container: container.InspectResponse{
			Config: &container.Config{
				Labels: map[string]string{
					"com.docker.compose.project.working_dir": "/services/service1",
				},
			},
		},
	}, nil)

	mockClient.On("ContainerInspect", mock.Anything, "container2", mock.Anything).Return(client.ContainerInspectResult{
		Container: container.InspectResponse{
			Config: &container.Config{
				Labels: map[string]string{
					"com.docker.compose.project.working_dir": "/services/service2",
				},
			},
		},
	}, nil)

	servicesDir := "/services"
	result, err := getManagedContainersWithClient(mockClient, servicesDir)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Contains(t, result, "service1")
	assert.Contains(t, result, "service2")
	assert.Len(t, result["service1"], 1)
	assert.Len(t, result["service2"], 1)

	// Test error case
	mockClient.On("ContainerList", mock.Anything, mock.Anything).Once().Return(client.ContainerListResult{}, errors.New("failed to list containers"))

	_, err = getManagedContainersWithClient(mockClient, servicesDir)

	assert.Error(t, err)
	assert.ErrorContains(t, err, "failed to list containers")
}

func TestGetServiceNameFromLabel(t *testing.T) {
	testCases := []struct {
		name           string
		containerID    string
		labels         map[string]string
		inspectError   error
		expectedResult string
		expectedError  error
	}{
		{
			name:        "Successful case",
			containerID: "container1",
			labels: map[string]string{
				"com.docker.compose.project.working_dir": "/services/service1",
			},
			inspectError:   nil,
			expectedResult: "service1",
			expectedError:  nil,
		},

		{
			name:           "Label not found",
			containerID:    "container2",
			labels:         map[string]string{},
			inspectError:   nil,
			expectedResult: "",
			expectedError:  nil,
		},

		{
			name:           "Error case",
			containerID:    "container3",
			labels:         map[string]string{},
			inspectError:   errors.New("failed to inspect container"),
			expectedResult: "",
			expectedError:  errors.New("failed to inspect container"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockClient := new(MockClient)

			mockClient.On("ContainerInspect", mock.Anything, tc.containerID, mock.Anything).Return(client.ContainerInspectResult{
				Container: container.InspectResponse{
					Config: &container.Config{
						Labels: tc.labels,
					},
				},
			}, tc.inspectError)

			servicesDir := "/services"
			serviceName, err := getServiceNameFromLabel(context.Background(), mockClient, container.Summary{ID: tc.containerID}, servicesDir)

			assert.Equal(t, tc.expectedResult, serviceName)
			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
