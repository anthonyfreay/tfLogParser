package cmd

import (
	"bytes"
	"errors"
	"flag"
	"log"
	"os"
	"testing"
)

// Mock implementation of the FilterFunc
func mockFilterLogs(filePath, level, startTime, endTime, searchKeyword string) error {
	// Simulate error for testing purposes
	if filePath == "error.log" {
		return errors.New("mock error")
	}
	return nil
}

func TestExecute(t *testing.T) {
	// Backup and restore the original command-line arguments
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()

	// Define test cases
	tests := []struct {
		name          string
		args          []string
		expectedError bool
		expectedLog   string
	}{
		{
			name:          "Missing log file",
			args:          []string{"cmd"},
			expectedError: true,
			expectedLog:   "Please provide a log file.\n",
		},
		{
			name:          "Valid log file with INFO level",
			args:          []string{"cmd", "test.log", "-level", "INFO"},
			expectedError: false,
			expectedLog:   "",
		},
		{
			name:          "Filter logs returns error",
			args:          []string{"cmd", "error.log", "-level", "ERROR"},
			expectedError: true,
			expectedLog:   "Error: mock error\n",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Override command-line arguments
			os.Args = tc.args

			// Reset the flag package (to avoid conflicts between tests)
			flag.CommandLine = flag.NewFlagSet(tc.name, flag.ExitOnError)

			// Capture output
			var buf bytes.Buffer
			log.SetOutput(&buf) // Redirect log output to buffer

			// Call Execute with the mock filter function
			Execute(mockFilterLogs)

			// Validate output based on expected error
			output := buf.String()
			if tc.expectedLog != output {
				t.Errorf("Expected log output %q but got %q", tc.expectedLog, output)
			}
		})
	}
}
