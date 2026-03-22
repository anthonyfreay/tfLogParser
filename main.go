package main

import (
	"fmt"
	"os"
	"tfLogParser/cmd"
	"tfLogParser/pkg/parser"
)

func main() {
	err := cmd.Execute(func(filePath, level, startTime, endTime, searchKeyword string) error {
		return parser.FilterLogsByLevelAndTimeAndKeyword(os.Stdout, filePath, level, startTime, endTime, searchKeyword)
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1) // Ensure the program exits with a non-zero status on error
	}
}
