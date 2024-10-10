package parser

import (
	"os"
	"testing"
)

func TestGetPriority(t *testing.T) {
	// Mocking logLevelPriority map for test
	logLevelPriority = map[string]int{
		"TRACE": 1,
		"DEBUG": 2,
		"INFO":  3,
		"WARN":  4,
		"ERROR": 5,
	}

	testCases := []struct {
		level    string
		expected int
	}{
		{level: "TRACE", expected: 1},   // Valid level
		{level: "DEBUG", expected: 2},   // Valid level
		{level: "INFO", expected: 3},    // Valid level
		{level: "WARN", expected: 4},    // Valid level
		{level: "ERROR", expected: 5},   // Valid level
		{level: "WARNING", expected: 0}, // Invalid level
		{level: "FOOBAR", expected: 0},  // Invalid level
		{level: "", expected: 0},        // Empty level (invalid)
	}

	for _, tc := range testCases {
		entry := LogEntry{Level: tc.level}
		result := entry.GetPriority()

		if result != tc.expected {
			t.Errorf("For level %s, expected priority %d, but got %d", tc.level, tc.expected, result)
		}
	}
}

func TestIsContinuationLine(t *testing.T) {
	testCases := []struct {
		logString                  string
		expectedIsContinuationLine bool
	}{
		{logString: "2024-10-03T00:43:29.930-0400 [TRACE] backend/local: state manager for workspace \"default\" will:", expectedIsContinuationLine: false},
		{logString: " - read initial snapshot from terraform.tfstate", expectedIsContinuationLine: true},
		{logString: " - write new snapshots to terraform.tfstate", expectedIsContinuationLine: true},
		{logString: " - create any backup at terraform.tfstate.backup", expectedIsContinuationLine: true},
		{logString: "2024-10-03T00:43:29.930-0400 [TRACE] backend/local: requesting state lock for workspace \"default\"", expectedIsContinuationLine: false},
	}

	for _, tc := range testCases {
		result := IsContinuationLine(tc.logString)

		if result != tc.expectedIsContinuationLine {
			t.Errorf("For log string %s, expected IsContinuationLine is: %t, but got %t", tc.logString, tc.expectedIsContinuationLine, result)
		}
	}
}

func TestParseLogLine(t *testing.T) {
	line := "2024-10-03T00:43:29.918-0400 [INFO] provider: Starting..."
	entry, err := ParseLogLine(line)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if entry.Level != "INFO" {
		t.Errorf("Expected level INFO, got %s", entry.Level)
	}
}

func TestIsWithinTimeRange(t *testing.T) {
	testCases := []struct {
		startTime           string
		endTime             string
		currentTime         string
		expectedInTimeRange bool
	}{
		{startTime: "2024-10-03T00:43:29.930-0400", endTime: "2024-10-03T00:43:29.931-0400", currentTime: "2024-10-03T00:43:29.932-0400", expectedInTimeRange: false},
		{startTime: "2024-10-03T00:43:29.930-0400", endTime: "2024-10-03T00:43:29.935-0400", currentTime: "2024-10-03T00:43:29.932-0400", expectedInTimeRange: true},
	}

	for _, tc := range testCases {
		result, err := IsWithinTimeRange(tc.currentTime, tc.startTime, tc.endTime)

		if err != nil {
			t.Error(err)
		}

		if result != tc.expectedInTimeRange {
			t.Errorf("For timerange %s to %s with current time %s, expected: %t, got: %t", tc.startTime, tc.endTime, tc.currentTime, tc.expectedInTimeRange, result)
		}
	}
}

func TestFilterLogsByLevelAndTimeAndKeyword(t *testing.T) {
	logContent := `
2024-10-03T00:43:29.918-0400 [INFO]  Terraform version: 1.9.5
2024-10-03T00:43:29.918-0400 [DEBUG] using github.com/hashicorp/go-tfe v1.58.0
2024-10-03T00:43:29.918-0400 [DEBUG] using github.com/hashicorp/hcl/v2 v2.20.0
2024-10-03T00:43:29.918-0400 [DEBUG] using github.com/hashicorp/terraform-svchost v0.1.1
2024-10-03T00:43:29.918-0400 [DEBUG] using github.com/zclconf/go-cty v1.14.4
`
	tmpfile, err := os.CreateTemp("", "testlog-*.log")
	if err != nil {
		t.Fatalf("unable to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(logContent)); err != nil {
		t.Fatalf("unable to write to temp file: %v", err)
	}

	if err := tmpfile.Close(); err != nil {
		t.Fatalf("unable to close temp file: %v", err)
	}

	startTime := "2024-10-03T00:00:00-0400"
	endTime := "2024-10-03T01:00:00-0400"

	testCases := []struct {
		name        string
		minLogLevel string
		startTime   string
		endTime     string
		keyword     string
		expectError bool
	}{
		{"Filter INFO Logs", "INFO", startTime, endTime, "", false},
		{"Filter DEBUG Logs and Keyword 'github'", "DEBUG", startTime, endTime, "github", false},
		{"Invalid Log Level", "INVALID", startTime, endTime, "", true},
		{"No Matching Time Range", "INFO", "2024-10-02T00:00:00-0400", "2024-10-02T23:59:59-0400", "", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := FilterLogsByLevelAndTimeAndKeyword(tmpfile.Name(), tc.minLogLevel, tc.startTime, tc.endTime, tc.keyword)
			if (err != nil) != tc.expectError {
				t.Errorf("unexpected error result: got %v, want error=%v", err, tc.expectError)
			}
		})
	}

}
