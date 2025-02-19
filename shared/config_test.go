package shared

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/pelletier/go-toml"
)

func TestSaveConfig(t *testing.T) {

	// Set a test configuration.
	now := time.Now().Truncate(time.Second)
	testConfig := ParetoConfig{
		TeamID:    "test-team",
		AuthToken: "test-token",
		Checks: map[string]CheckStatus{
			"check1": {
				UpdatedAt: now,
				Passed:    true,
				Disabled:  false,
			},
		},
	}
	Config = testConfig

	// Call SaveConfig.
	if err := SaveConfig(); err != nil {
		t.Fatalf("SaveConfig failed: %v", err)
	}
	configDir, err := os.UserConfigDir()
	if err != nil {
		t.Fatalf("Failed to read the saved config file: %v", err)
	}

	// The configuration file should be located in the temporary config directory.
	configPath := filepath.Join(configDir, "pareto.toml")
	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("Failed to read the saved config file: %v", err)
	}

	var loadedConfig ParetoConfig
	if err := toml.Unmarshal(data, &loadedConfig); err != nil {
		t.Fatalf("Failed to unmarshal config data: %v", err)
	}

	// Check that the saved configuration matches the test configuration.
	if loadedConfig.TeamID != testConfig.TeamID {
		t.Errorf("Expected TeamID %q, got %q", testConfig.TeamID, loadedConfig.TeamID)
	}
	if loadedConfig.AuthToken != testConfig.AuthToken {
		t.Errorf("Expected AuthToken %q, got %q", testConfig.AuthToken, loadedConfig.AuthToken)
	}
	// Since time values can be sensitive, compare using Truncate.
	if len(loadedConfig.Checks) != len(testConfig.Checks) {
		t.Errorf("Expected %d checks, got %d", len(testConfig.Checks), len(loadedConfig.Checks))
	}

	for key, expected := range testConfig.Checks {
		actual, ok := loadedConfig.Checks[key]
		if !ok {
			t.Errorf("Expected check %q not found", key)
			continue
		}
		// Compare time values after truncation.
		if !actual.UpdatedAt.Truncate(time.Second).Equal(expected.UpdatedAt) {
			t.Errorf("For check %q, expected UpdatedAt %v, got %v", key, expected.UpdatedAt, actual.UpdatedAt)
		}
		if actual.Passed != expected.Passed {
			t.Errorf("For check %q, expected Passed %v, got %v", key, expected.Passed, actual.Passed)
		}
		if actual.Disabled != expected.Disabled {
			t.Errorf("For check %q, expected Disabled %v, got %v", key, expected.Disabled, actual.Disabled)
		}
	}
}
