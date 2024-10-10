package cmd

import (
	"bytes"
	"errors"
	"flag"
	"log"
	"os"
	"strings"
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
			args:          []string{"cmd", "--level=INFO", "--start-time=2024-10-03T00:43:29.918-0400", "--end-time=2024-10-03T00:43:29.919-0400"},
			expectedError: true,
			expectedLog:   "please provide a log file path using the -file flag\n",
		},
		{
			name:          "Valid log file with INFO level",
			args:          []string{"cmd", "--level=INFO", "--start-time=2024-10-03T00:43:29.918-0400", "--end-time=2024-10-03T00:43:29.919-0400", "--file=resources/log.txt"},
			expectedError: false,
			expectedLog:   "",
		},
		{
			name:          "Filter logs returns error",
			args:          []string{"cmd", "--level=ERROR", "--start-time=2024-10-03T00:43:29.918-0400", "--end-time=2024-10-03T00:43:29.919-0400", "--file=error.log"},
			expectedError: true,
			expectedLog:   "error filtering logs: mock error\n",
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
			err := Execute(mockFilterLogs)

			// Check if an error was expected
			if (err != nil) != tc.expectedError {
				t.Errorf("Expected error: %v, got: %v", tc.expectedError, err != nil)
			}

			// Validate log output using strings.Contains
			output := buf.String()
			if tc.expectedLog != "" && !strings.Contains(output, tc.expectedLog) {
				t.Errorf("Expected log to contain %q but got %q", tc.expectedLog, output)
			}
		})
	}
}
