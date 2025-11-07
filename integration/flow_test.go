package integration

import (
	"context"
	"embed"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/docker/compose/v2/pkg/api"
	"github.com/testcontainers/testcontainers-go/modules/compose"
)

//go:embed test_data/*
var testDataFS embed.FS

func TestFileGeneration(t *testing.T) {
	ctx := context.Background()

	baseDir := t.TempDir()
	//baseDir := t.TempDir()
	servicesDir := filepath.Join(baseDir, "services")
	dataDir := filepath.Join(baseDir, "data")
	configDir := filepath.Join(baseDir, "config")
	os.Mkdir(servicesDir, 0o777)
	os.Mkdir(dataDir, 0o777)
	os.Mkdir(configDir, 0o777)
	err := os.WriteFile(filepath.Join(configDir, "config.yaml"),
		[]byte(strings.Join([]string{

			"AUTONAS_HOST: truenas-scale",
			"SERVICES_PATH: " + servicesDir,
			"DATA_PATH: " + dataDir,
			"enabled_services:",
			"  - homepage",
		}, "\n")), 0o777)
	if err != nil {
		t.Fatalf(" error while creating config file  : %v\n", err)
	}
	const configFiles = "config.default.yaml,/config/config.yaml"

	//server, _ := testutil.NewGitTestServer()
	//defer server.Close()

	// Add test files
	//server.AddFile("README.md", "# Test Repo")

	// Start docker-compose environment
	composeEnv, err := compose.NewDockerCompose("../compose.yaml")
	composeEnv.WithEnv(map[string]string{
		"CONFIG_FILES":  configFiles,
		"CONFIG_REPO":   "https://github.com/omar-kada/autonas-config",
		"CRON_PERIOD":   "*/10 * * * *",
		"SERVICES_PATH": servicesDir,
		"CONFIG_PATH":   configDir,
		"ENV":           "DEV",
		"UID":           fmt.Sprint(os.Getuid()),
		"GID":           fmt.Sprint(os.Getgid()),
	})
	if err != nil {
		t.Fatalf("failed to load compose: %v", err)
	}

	// Start containers
	if err := composeEnv.Up(ctx, compose.WithRecreate(api.RecreateForce)); err != nil {
		t.Fatalf("failed to start compose: %v", err)
	}

	// Ensure environment is stopped after test
	t.Cleanup(func() {
		composeEnv.Down(ctx, compose.RemoveOrphans(true), compose.RemoveVolumes(true))
	})

	// Wait for file to be generated (polling)
	targetFile := filepath.Join(servicesDir, "homepage", ".env")

	ok := waitForFile(targetFile, 10*time.Minute)
	if !ok {

		cont, err := composeEnv.ServiceContainer(ctx, "autonas")
		if err != nil {
			t.Fatal(err)
		}
		logsReader, err := cont.Logs(ctx)
		if err != nil {
			t.Fatal(err)
		}
		state, err := cont.State(ctx)
		if err != nil {
			t.Fatal(err)
		}

		bytes, _ := io.ReadAll(logsReader)
		fmt.Printf("Logs for service %s:\n%s\n%v\n", "autonas", string(bytes), state)
		t.Errorf("expected file was not generated: %s", targetFile)
	}
}

func waitForFile(path string, timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		if _, err := os.Stat(path); err == nil {
			return true
		}
		time.Sleep(200 * time.Millisecond)
	}
	return false
}

func LocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}

	for _, addr := range addrs {
		// Check the address type and make sure it's not a loopback
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ip4 := ipNet.IP.To4(); ip4 != nil {
				return ip4.String()
			}
		}
	}

	return ""
}
