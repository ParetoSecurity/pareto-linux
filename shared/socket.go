package shared

import (
	"encoding/json"
	"net"
	"sync"
	"time"

	"github.com/caarlos0/log"
	"github.com/davecgh/go-spew/spew"
)

var (
	rateLimit sync.Mutex
	lastCall  time.Time
)

func RunCheckViaHelper(uuid string) (bool, error) {
	rateLimit.Lock()
	defer rateLimit.Unlock()

	for time.Since(lastCall) < time.Second*2 {
		time.Sleep(time.Millisecond * 100)
	}
	lastCall = time.Now()

	conn, err := net.Dial("unix", "/run/pareto.sock")
	if err != nil {
		log.WithError(err).Warn("Failed to connect to helper")
		return false, err
	}
	defer conn.Close()

	// Send UUID
	input := map[string]string{"uuid": uuid}
	encoder := json.NewEncoder(conn)
	log.WithField("input", spew.Sdump(input)).Debug("Sending input to helper")
	if err := encoder.Encode(input); err != nil {
		log.WithError(err).Warn("Failed to encode JSON")
		return false, err
	}

	// Read response
	decoder := json.NewDecoder(conn)
	var status map[string]bool
	if err := decoder.Decode(&status); err != nil {
		log.WithError(err).Warn("Failed to decode JSON")
		return false, err
	}
	log.WithField("status", status).Debug("Received status from helper")
	return status[uuid], nil
}
