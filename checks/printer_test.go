package checks

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrinterRun(t *testing.T) {
	tests := []struct {
		name           string
		mockCheckPort  func(port int, proto string) bool
		expectedPassed bool
		expectedPorts  map[int]string
	}{
		{
			name: "No ports open",
			mockCheckPort: func(port int, proto string) bool {
				return false
			},
			expectedPassed: true,
			expectedPorts:  map[int]string{},
		},
		{name: "CUPS port open", mockCheckPort: func(port int, proto string) bool {
			return port == 631
		}, expectedPassed: false, expectedPorts: map[int]string{
			631: "CUPS",
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			checkPortMock = tt.mockCheckPort
			printer := &Printer{}

			err := printer.Run()
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedPassed, printer.Passed())
			assert.Equal(t, tt.expectedPorts, printer.ports)
		})
	}
}
