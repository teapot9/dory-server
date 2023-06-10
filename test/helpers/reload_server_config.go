package helpers

import (
	"encoding/json"
	"os"
	"testing"
	"time"
	"github.com/be-ys-cloud/dory-server/internal/structures"
)

// ReloadServerConfig reload configuration.json and restart the server container
func ReloadServerConfig(t *testing.T, cfg *structures.Configuration, ct *ContainersEnvironment) {
	cfgBin, err := json.Marshal(cfg)
	if err != nil {
		t.Fatalf("failed to marshal new config: %v", err)
	}

	path, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}

	err = os.WriteFile(path + "/configuration.json", cfgBin, 0777)
	if err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	err = ct.Pool.Client.RestartContainer(ct.Server.Container.ID, 10)
	if err != nil {
		t.Fatalf("failed to restart server container: %v", err)
	}

	time.Sleep(1 * time.Second)
}
