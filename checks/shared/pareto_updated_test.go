package shared

import (
	"testing"
)

func TestParetoUpdated_Run(t *testing.T) {

	check := &ParetoUpdated{}
	err := check.Run()

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

}
