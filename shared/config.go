package shared

import (
	"os"
	"path/filepath"
	"time"

	"github.com/pelletier/go-toml"
)

var Config ParetoConfig

type CheckStatus struct {
	UpdatedAt time.Time
	Passed    bool
	Disabled  bool
}

type ParetoConfig struct {
	TeamID    string
	AuthToken string
	Checks    map[string]CheckStatus
}

func SaveConfig() error {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return err
	}

	configPath := filepath.Join(configDir, "pareto.toml")
	file, err := os.Create(configPath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := toml.NewEncoder(file)
	return encoder.Encode(Config)
}

func LoadConfig() error {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return err
	}

	configPath := filepath.Join(configDir, "pareto.toml")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		if err := SaveConfig(); err != nil {
			return err
		}
		return nil
	}
	file, err := os.Open(configPath)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := toml.NewDecoder(file)
	err = decoder.Decode(&Config)
	if err != nil {
		return err
	}

	return nil
}
