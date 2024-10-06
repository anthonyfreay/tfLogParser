package cmd

import (
	"flag"
	"fmt"
	"tfLogParser/pkg/parser"
)

func Execute() {
	logLevel := flag.String("level", "INFO", "Minimum log level to display (TRACE, DEBUG, INFO, WARN, ERROR)")
	startTime := flag.String("start-time", "", "Filter logs starting from this time (RFC3339 format)")
	endTime := flag.String("end-time", "", "Filter logs up to this time (RFC3339 format)")
	flag.Parse()

	if len(flag.Args()) == 0 {
		fmt.Println("Please provide a log file.")
		return
	}

	filePath := flag.Arg(0)

	// Call the log parser and filter by log level and time range
	err := parser.FilterLogsByLevelAndTime(filePath, *logLevel, *startTime, *endTime)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
