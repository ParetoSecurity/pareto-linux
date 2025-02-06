package claims

import (
	"testing"
)

func TestClaims(t *testing.T) {
	for _, claim := range All {
		for _, check := range claim.Checks {
			if check == nil {
				t.Errorf("Claim %s has a nil check", claim.Title)
			}
			check.RequiresRoot()
			check.Status()
			check.UUID()
			check.Passed()
		}
	}
}
