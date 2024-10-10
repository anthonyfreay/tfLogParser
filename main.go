package main

import (
	"tfLogParser/cmd"
	"tfLogParser/pkg/parser"
)

func main() {
	cmd.Execute(parser.FilterLogsByLevelAndTimeAndKeyword)
}
