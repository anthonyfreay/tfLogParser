package cmd

import (
	"flag"
	"fmt"
)

// Type definition for the filter function
type FilterFunc func(filePath, level, startTime, endTime, searchKeyword string) error

// Execute accepts the filtering function and file path as a dedicated flag
func Execute(filterLogs FilterFunc) error {
	logLevel := flag.String("level", "INFO", "Minimum log level to display (TRACE, DEBUG, INFO, WARN, ERROR)")
	startTime := flag.String("start-time", "", "Filter logs starting from this time (RFC3339 format)")
	endTime := flag.String("end-time", "", "Filter logs up to this time (RFC3339 format)")
	searchKeyword := flag.String("search", "", "Keyword to search for in log messages")
	filePath := flag.String("file", "", "Path to the log file to be processed")

	flag.Parse()

	if *filePath == "" {
		return fmt.Errorf("please provide a log file path using the -file flag")
	}

	// Call the injected filter function with the provided file path and other flags
	err := filterLogs(*filePath, *logLevel, *startTime, *endTime, *searchKeyword)
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	return nil
}
