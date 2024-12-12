package shared

import (
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/caarlos0/log"
	"github.com/pelletier/go-toml"
)

type LastState struct {
	UUID    string `json:"uuid"`
	State   bool   `json:"state"`
	Details string `json:"details"`
}

var (
	mutex       sync.Mutex
	states      map[string]LastState
	lastModTime time.Time
	configPath  string
)

func init() {
	states = make(map[string]LastState)
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.WithError(err).Fatal("failed to get user home directory")
	}
	configPath = filepath.Join(homeDir, ".cache", "paretosecurity.state")
}

// Commit writes the current state map to the TOML file.
func CommitLastState() error {
	file, err := os.Create(configPath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := toml.NewEncoder(file)
	return encoder.Encode(states)
}

// UpdateState updates the LastState struct in the in-memory map and commits to the TOML file.
func UpdateLastState(newState LastState) {
	mutex.Lock()
	defer mutex.Unlock()

	states[newState.UUID] = newState
}

// GetState retrieves the LastState struct by UUID.
func GetLastState(uuid string) (LastState, bool, error) {
	mutex.Lock()
	defer mutex.Unlock()

	fileInfo, err := os.Stat(configPath)
	if err != nil {
		return LastState{}, false, err
	}

	if fileInfo.ModTime().After(lastModTime) {
		file, err := os.Open(configPath)
		if err != nil {
			return LastState{}, false, err
		}
		defer file.Close()

		decoder := toml.NewDecoder(file)
		if err := decoder.Decode(&states); err != nil {
			return LastState{}, false, err
		}

		lastModTime = fileInfo.ModTime()
	}

	state, exists := states[uuid]
	return state, exists, nil
}
