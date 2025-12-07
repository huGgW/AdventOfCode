package main

import (
	"bufio"
	"day4/grid"
	"flag"
	"fmt"
	"io"
	"os"
)

var (
	inputFileFlag string
)

func main() {
	flag.StringVar(&inputFileFlag, "input", "", "Path to input file")
	flag.Parse()
	if inputFileFlag == "" {
		panic("input file flag is required")
	}

	f, err := os.Open(inputFileFlag)
	if err != nil {
		panic(fmt.Sprintf("failed to open input file: %w", err))
	}

	gr, err := ParseInput(f)
	if err != nil {
		panic(fmt.Sprintf("failed to parse grid: %w", err))
	}

	accessables := gr.ForkLiftAccessablePosition()
	fmt.Printf("%d rolls of paper can be accessed by a forklift.\n", len(accessables))

	var maxRemoved int
	for len(accessables) > 0 {
		maxRemoved += len(accessables)

		if err := gr.RemovePapers(accessables); err != nil {
			panic(fmt.Sprintf("failed to remove papers: %w", err))
		}

		accessables = gr.ForkLiftAccessablePosition()
	}

	fmt.Printf("Maximum of %d rolls of paper can be removed by forklift.\n", maxRemoved)
}

func ParseInput(r io.Reader) (grid.Grid, error) {
	gb := grid.GridBuilder{}
	sc := bufio.NewScanner(r)

	for {
		scanned := sc.Scan()
		if !scanned {
			if err := sc.Err(); err != nil {
				return grid.Grid{}, fmt.Errorf("error scanning input: %w", err)
			}
			break
		}

		line := sc.Text()
		runes := []rune(line)
		gb.Append(runes)
	}

	gr, err := gb.Build()
	if err != nil {
		return grid.Grid{}, fmt.Errorf("failed to build grid: %w", err)
	}

	return gr, nil
}
