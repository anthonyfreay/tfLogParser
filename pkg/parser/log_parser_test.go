package parser

import "testing"

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
