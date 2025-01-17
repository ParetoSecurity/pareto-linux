package claims

import (
	"testing"
)

func TestClaims(t *testing.T) {
	for _, claim := range All {
		if len(claim.Checks) == 0 {
			t.Errorf("Claim %s has no checks", claim.Title)
		}
		for _, check := range claim.Checks {
			if check == nil {
				t.Errorf("Claim %s has a nil check", claim.Title)
			}
			check.RequiresRoot()
			check.ReportIfDisabled()
			check.Status()
			check.UUID()
			check.Passed()
		}
	}
}
