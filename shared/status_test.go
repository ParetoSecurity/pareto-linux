package shared

import "testing"

func TestIsLinked(t *testing.T) {

	tests := []struct {
		name      string
		teamID    string
		authToken string
		expected  bool
	}{
		{"both empty", "", "", false},
		{"only teamID set", "team123", "", false},
		{"only auth token set", "", "token123", false},
		{"both set", "team123", "token123", true},
	}

	for _, tt := range tests {
		Config.TeamID = tt.teamID
		Config.AuthToken = tt.authToken
		if got := IsLinked(); got != tt.expected {
			t.Errorf("%s: expected %v, got %v", tt.name, tt.expected, got)
		}
	}
}
