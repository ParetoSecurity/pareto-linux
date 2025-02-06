package checks

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSecureBoot_Run(t *testing.T) {
	tests := []struct {
		name           string
		mockFiles      map[string][]byte
		expectedPassed bool
		expectedStatus string
	}{
		{
			name: "SecureBoot enabled",
			mockFiles: map[string][]byte{
				"/sys/firmware/efi/efivars/SecureBoot-1234": {0, 0, 0, 0, 1},
			},
			expectedPassed: true,
			expectedStatus: "SecureBoot is enabled",
		},
		{
			name: "SecureBoot disabled",
			mockFiles: map[string][]byte{
				"/sys/firmware/efi/efivars/SecureBoot-1234": {0, 0, 0, 0, 0},
			},
			expectedPassed: false,
			expectedStatus: "SecureBoot is disabled",
		},
		{
			name:           "SecureBoot EFI variable not found",
			mockFiles:      map[string][]byte{},
			expectedPassed: false,
			expectedStatus: "Could not find SecureBoot EFI variable",
		},
		{
			name: "SecureBoot EFI variable read error",
			mockFiles: map[string][]byte{
				"/sys/firmware/efi/efivars/SecureBoot-1234": nil,
			},
			expectedPassed: false,
			expectedStatus: "Could not read SecureBoot status",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock os.ReadFile
			osReadFile = func(name string) ([]byte, error) {
				if data, ok := tt.mockFiles[name]; ok {
					return data, nil
				}
				return nil, os.ErrNotExist
			}

			sb := &SecureBoot{}
			err := sb.Run()
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedPassed, sb.Passed())
			assert.Equal(t, tt.expectedStatus, sb.Status())
		})
	}
}

func TestSecureBoot_IsRunnable(t *testing.T) {
	tests := []struct {
		name           string
		mockStatError  error
		expectedResult bool
		expectedStatus string
	}{
		{
			name:           "System running in UEFI mode",
			mockStatError:  nil,
			expectedResult: false,
			expectedStatus: "",
		},
		{
			name:           "System not running in UEFI mode",
			mockStatError:  os.ErrNotExist,
			expectedResult: true,
			expectedStatus: "System is not running in UEFI mode",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock os.Stat
			osStat = func(name string) (os.FileInfo, error) {
				return nil, tt.mockStatError
			}

			sb := &SecureBoot{}
			result := sb.IsRunnable()
			assert.Equal(t, tt.expectedResult, result)
			assert.Equal(t, tt.expectedStatus, sb.Status())
		})
	}
}
