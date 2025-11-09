package integration

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/docker/compose/v2/pkg/api"
	"github.com/docker/docker/api/types/container"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/compose"
)

func TestFileGeneration(t *testing.T) {
	ctx := context.Background()

	// Given

	baseDir := t.TempDir()
	if os.Getenv("GITHUB_RUNNER") == "true" {
		// disable auto clean up of temp directory in github runner (because of permissions issue)
		// it will be handled by the runner itsef
		baseDir = "/tmp/integration-tests"
		os.MkdirAll(baseDir, 0777)
	}

	servicesDir := filepath.Join(baseDir, "services")
	dataDir := filepath.Join(baseDir, "data")
	configDir := filepath.Join(baseDir, "config")
	os.Mkdir(servicesDir, 0777)
	os.Mkdir(dataDir, 0777)
	os.Mkdir(configDir, 0777)

	os.Chmod(servicesDir, 0777)
	os.Chmod(dataDir, 0777)

	err := os.WriteFile(filepath.Join(configDir, "config.yaml"),
		[]byte(strings.Join([]string{
			"AUTONAS_HOST: test",
			"SERVICES_PATH: " + servicesDir,
			"DATA_PATH: " + dataDir,
			"enabled_services:",
			"  - homepage",
		}, "\n")), 0777)
	assert.NoError(t, err, "error while creating config file")

	const configFiles = "config.default.yaml,/config/config.yaml"

	// When
	// Start docker-compose environment
	composeEnv, err := compose.NewDockerCompose("../compose.yaml")
	composeEnv.WithEnv(map[string]string{
		"BUILD_CONTEXT": ".",
		"VARSION":       "local",
		"CONFIG_FILES":  configFiles,
		"CONFIG_REPO":   "https://github.com/omar-kada/autonas-config",
		"CRON_PERIOD":   "*/10 * * * *",
		"SERVICES_PATH": servicesDir,
		"CONFIG_PATH":   configDir,
		"ENV":           "DEV",
		"UID":           fmt.Sprint(os.Getuid()),
		"GID":           fmt.Sprint(os.Getgid()),
	})

	assert.NoError(t, err, "failed to load compose")

	// Start containers
	if err := composeEnv.Up(ctx, compose.WithRecreate(api.RecreateForce)); err != nil {
		t.Fatalf("failed to start compose: %v", err)
	}
	t.Cleanup(func() {
		if t.Failed() {
			printContainerLogs(ctx, t, composeEnv)
		}
		composeEnv.Down(ctx, compose.RemoveOrphans(true), compose.RemoveVolumes(true))
	})

	// Then
	// Wait for file to be generated (polling)
	targetFile := filepath.Join(servicesDir, "homepage", ".env")

	ok := waitForFile(targetFile, 1*time.Minute)
	if !ok {
		t.Errorf("expected file was not generated: %s", targetFile)
	}
	targetWorkDir := filepath.Join(servicesDir, "homepage")

	homepageContainers, _ := waitForComposeStack(ctx, targetWorkDir, 2*time.Minute)
	assert.NotEmpty(t, homepageContainers, "homepage container not found")
}

func printContainerLogs(ctx context.Context, t *testing.T, composeEnv *compose.DockerCompose) {
	t.Helper()

	for _, service := range composeEnv.Services() {

		cont, err := composeEnv.ServiceContainer(ctx, service)
		assert.NoError(t, err)

		logsReader, err := cont.Logs(ctx)
		assert.NoError(t, err)

		bytes, err := io.ReadAll(logsReader)

		assert.NoError(t, err)
		fmt.Printf("Logs for service %s:\n%s\n", service, string(bytes))
	}
}

func waitForFile(path string, timeout time.Duration) bool {
	res, _ := waitFor(timeout, func() (bool, bool) {
		_, err := os.Stat(path)
		return err == nil, err == nil
	})
	return res
}

func waitForComposeStack(ctx context.Context, workingDir string, timeout time.Duration) ([]*container.Summary, bool) {
	container, ok := waitFor(timeout, func() ([]*container.Summary, bool) {
		containers, err := findComposeStack(ctx, workingDir)
		return containers, err == nil && len(containers) > 0
	})
	return container, ok
}

func waitFor[T any](timeout time.Duration, predicate func() (T, bool)) (T, bool) {
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		if res, ok := predicate(); ok {
			return res, true
		}
		time.Sleep(200 * time.Millisecond)
	}
	var zero T
	return zero, false
}

func findComposeStack(ctx context.Context, workingDir string) ([]*container.Summary, error) {
	// Create a Docker provider
	provider, err := testcontainers.NewDockerProvider()
	if err != nil {
		return nil, err
	}

	// List all running containers
	containers, err := provider.Client().ContainerList(ctx, container.ListOptions{All: true})
	if err != nil {
		return nil, err
	}
	workDirLabel := "com.docker.compose.project.working_dir"
	stackContainers := make([]*container.Summary, 0)
	// Look for Git server containers
	for _, container := range containers {
		if container.Labels[workDirLabel] == workingDir {
			stackContainers = append(stackContainers, &container)
		}
	}

	return stackContainers, nil
}
