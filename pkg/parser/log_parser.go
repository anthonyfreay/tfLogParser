package parser

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"time"
)

// LogEntry represents a parsed log entry.
type LogEntry struct {
	Timestamp string
	Level     string
	Component string
	Message   string
}

// logLevelPriority maps log levels to numeric values for filtering.
var logLevelPriority = map[string]int{
	"TRACE": 1,
	"DEBUG": 2,
	"INFO":  3,
	"WARN":  4,
	"ERROR": 5,
}

// GetPriority returns the numeric priority of a log level.
func (e *LogEntry) GetPriority() int {
	if priority, exists := logLevelPriority[e.Level]; exists {
		return priority
	}
	return 0 // Default for unknown levels.
}

var (
	timestampPattern = `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}`
	timestampRegex   = regexp.MustCompile(timestampPattern)
	logPattern       = `(?P<Timestamp>\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}.\d+-\d{4}) \[(?P<Level>\w+)\]\s+(?P<Component>[\w\/\.\s]*)?:?\s*(?P<Message>.+)`
	logRegex         = regexp.MustCompile(logPattern)
)

// IsContinuationLine checks if a line is a continuation (i.e., has no timestamp or log level).
func IsContinuationLine(line string) bool {
	return !timestampRegex.MatchString(line)
}

// ParseLogLine parses a single log line into a LogEntry struct.
func ParseLogLine(line string) (*LogEntry, error) {
	matches := logRegex.FindStringSubmatch(line)

	if len(matches) == 0 {
		return nil, fmt.Errorf("could not parse log line: %s", line)
	}

	return &LogEntry{
		Timestamp: matches[1],
		Level:     matches[2],
		Component: strings.TrimSpace(matches[3]), // Trim extra spaces in the component.
		Message:   matches[4],
	}, nil
}

const layout = "2006-01-02T15:04:05-0700"

// IsWithinTimeRange checks if the given timestamp is within the start and end time range.
func IsWithinTimeRange(timestamp string, start, end *time.Time) (bool, error) {
	logTime, err := time.Parse(layout, timestamp)
	if err != nil {
		return false, fmt.Errorf("error parsing log timestamp: %w", err)
	}

	if start != nil && logTime.Before(*start) {
		return false, nil
	}

	if end != nil && logTime.After(*end) {
		return false, nil
	}

	return true, nil
}

// printEntry prints a LogEntry to the provided writer if it passes all filters
func printEntry(w io.Writer, entry *LogEntry, minLogLevelPriority int, start, end *time.Time, keyword string) error {
	if entry == nil {
		return nil
	}

	// Time filtering
	result, err := IsWithinTimeRange(entry.Timestamp, start, end)
	if err != nil {
		return fmt.Errorf("error checking time range: %v", err)
	}
	if !result {
		return nil
	}

	// Log level filtering
	if entry.GetPriority() < minLogLevelPriority {
		return nil
	}

	// Keyword filtering
	if keyword != "" && !strings.Contains(entry.Message, keyword) {
		return nil
	}

	// Print valid entry
	fmt.Fprintf(w, "%s [%s] %s: %s\n", entry.Timestamp, entry.Level, entry.Component, entry.Message)
	return nil
}

// FilterLogsByLevelAndTimeAndKeyword filters logs by level, time range, and keyword
func FilterLogsByLevelAndTimeAndKeyword(w io.Writer, filePath string, minLogLevel string, startTimeStr string, endTimeStr string, keyword string) error {
	minLogLevelPriority, exists := logLevelPriority[strings.ToUpper(minLogLevel)]
	if !exists {
		return fmt.Errorf("invalid log level: %s", minLogLevel)
	}

	var start, end *time.Time
	if startTimeStr != "" {
		t, err := time.Parse(layout, startTimeStr)
		if err != nil {
			return fmt.Errorf("error parsing startTime: %w", err)
		}
		start = &t
	}
	if endTimeStr != "" {
		t, err := time.Parse(layout, endTimeStr)
		if err != nil {
			return fmt.Errorf("error parsing endTime: %w", err)
		}
		end = &t
	}

	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	var currentEntry *LogEntry

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		if IsContinuationLine(line) {
			if currentEntry != nil {
				currentEntry.Message += " " + strings.TrimSpace(line)
			}
			continue
		}

		// Before parsing the new line, print the previous entry if it exists
		if currentEntry != nil {
			if err := printEntry(w, currentEntry, minLogLevelPriority, start, end, keyword); err != nil {
				return err
			}
		}

		entry, err := ParseLogLine(line)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing log line: %v\n", err)
			currentEntry = nil
			continue
		}

		currentEntry = entry
	}

	// Print the last entry after the loop finishes
	if currentEntry != nil {
		if err := printEntry(w, currentEntry, minLogLevelPriority, start, end, keyword); err != nil {
			return err
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading log file: %v", err)
	}

	return nil
}
