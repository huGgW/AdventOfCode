package util

import "flag"

func InitFlags() string {
	inputPath := flag.String(
		"input",
		"",
		"Path to input file",
	)

	flag.Parse()

	if inputPath == nil || *inputPath == "" {
		panic("input path is required")
	}

	return *inputPath
}
