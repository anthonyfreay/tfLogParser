package parser

import (
	"bufio"
	"fmt"
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

// IsContinuationLine checks if a line is a continuation (i.e., has no timestamp or log level).
func IsContinuationLine(line string) bool {
	timestampPattern := `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}`
	matched, _ := regexp.MatchString(timestampPattern, line)
	return !matched // Continuation line has no timestamp
}

// ParseLogLine parses a single log line into a LogEntry struct.
func ParseLogLine(line string) (*LogEntry, error) {
	logPattern := `(?P<Timestamp>\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}.\d+-\d{4}) \[(?P<Level>\w+)\]\s+(?P<Component>[\w\/\.\s]*)?:?\s*(?P<Message>.+)`
	re := regexp.MustCompile(logPattern)
	matches := re.FindStringSubmatch(line)

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

// IsWithinTimeRange checks if the given timestamp is within the startTime and endTime range
func IsWithinTimeRange(timestamp string, startTime string, endTime string) (bool, error) {
	layout := "2006-01-02T15:04:05-0700"

	logTime, err := time.Parse(layout, timestamp)
	if err != nil {
		return false, fmt.Errorf("error parsing log timestamp: %w", err)
	}

	// If both startTime and endTime are empty, allow all logs
	if startTime == "" && endTime == "" {
		return true, nil
	}

	// If no startTime is provided, treat it as no lower bound (i.e., allow all logs before endTime)
	if startTime != "" {
		start, err := time.Parse(layout, startTime)
		if err != nil {
			return false, fmt.Errorf("error parsing startTime: %w", err)
		}
		if logTime.Before(start) {
			return false, nil // Skip log if it is before the start time
		}
	}

	// If no endTime is provided, treat it as no upper bound (i.e., allow all logs after startTime)
	if endTime != "" {
		end, err := time.Parse(layout, endTime)
		if err != nil {
			return false, fmt.Errorf("error parsing endTime: %w", err)
		}
		if logTime.After(end) {
			return false, nil // Skip log if it is after the end time
		}
	}

	return true, nil // Allow log if it passes all checks
}

// FilterLogsByLevelAndKeyword filters logs by level, time range, and keyword
func FilterLogsByLevelAndTimeAndKeyword(filePath string, minLogLevel string, startTime string, endTime string, keyword string) error {
	minLogLevelPriority, exists := logLevelPriority[strings.ToUpper(minLogLevel)]
	if !exists {
		return fmt.Errorf("invalid log level: %s", minLogLevel)
	}

	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	var lastEntry *LogEntry

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Handle multi-line continuation
		if IsContinuationLine(line) {
			if lastEntry != nil {
				lastEntry.Message += " " + strings.TrimSpace(line)
			}
			continue
		}

		entry, err := ParseLogLine(line)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing log line: %v\n", err)
			continue
		}

		lastEntry = entry

		// Time filtering
		result, err := IsWithinTimeRange(entry.Timestamp, startTime, endTime)
		if err != nil {
			return fmt.Errorf("error checking time range: %v", err)
		}

		if !result {
			continue
		}

		// Log level filtering
		if entry.GetPriority() < minLogLevelPriority {
			continue
		}

		// Keyword filtering: if a keyword is provided, check if it exists in the log message
		if keyword != "" && !strings.Contains(entry.Message, keyword) {
			continue
		}

		// If log passes all filters, print it
		fmt.Printf("%s [%s] %s: %s\n", entry.Timestamp, entry.Level, entry.Component, entry.Message)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading log file: %v", err)
	}

	return nil
}
