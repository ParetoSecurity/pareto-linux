package shared

import (
	"encoding/json"
	"net"
)

func RunCheckViaHelper(uuid string) (bool, error) {
	conn, err := net.Dial("unix", "/var/run/pareto-linux.sock")
	if err != nil {
		return false, err
	}
	defer conn.Close()

	// Send UUID
	input := map[string]string{"uuid": uuid}
	encoder := json.NewEncoder(conn)
	if err := encoder.Encode(input); err != nil {
		return false, err
	}

	// Read response
	decoder := json.NewDecoder(conn)
	var status map[string]bool
	if err := decoder.Decode(&status); err != nil {
		return false, err
	}

	return status[uuid], nil
}
