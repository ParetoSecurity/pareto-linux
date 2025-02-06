package checks

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSharing_Run(t *testing.T) {
	tests := []struct {
		name     string
		mockFunc func(port int, proto string) bool
		expected bool
	}{
		{
			name: "No ports open",
			mockFunc: func(port int, proto string) bool {
				return false
			},
			expected: true,
		},
		{
			name: "Some ports open",
			mockFunc: func(port int, proto string) bool {
				if port == 445 || port == 2049 {
					return true
				}
				return false
			},
			expected: false,
		},
		{
			name: "All ports open",
			mockFunc: func(port int, proto string) bool {
				return true
			},
			expected: false,
		},
		{
			name: "Only one port open",
			mockFunc: func(port int, proto string) bool {
				if port == 445 {
					return true
				}
				return false
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock the checkPort function
			checkPortMock = tt.mockFunc

			sharing := &Sharing{}
			err := sharing.Run()
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, sharing.Passed())
			assert.NotEmpty(t, sharing.UUID())
			assert.False(t, sharing.RequiresRoot())
		})
	}
}
