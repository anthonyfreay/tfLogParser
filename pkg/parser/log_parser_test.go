package parser

import "testing"

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
	// TODO
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
	// TODO
}

func TestFilterLogsByLevelAndTimeAndKeyword(t *testing.T) {
	// TODO
}
