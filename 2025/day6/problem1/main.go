package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"day6/cephalopod"
	"day6/operator"
	"day6/util"
)

var whiteSpaceRegex *regexp.Regexp

func main() {
	inputPath := util.InitFlags()
	initRegex()

	f, err := os.Open(inputPath)
	if err != nil {
		log.Fatalf("failed to open input file: %v", err)
	}
	defer f.Close()

	cephalopodRows, err := parseInput(f)
	if err != nil {
		log.Fatalf("failed to parse input: %v", err)
	}

	calculatedResults := make([]int64, 0, len(cephalopodRows))
	for _, cephalopodRow := range cephalopodRows {
		calculatedResults = append(calculatedResults, cephalopodRow.Calculate())
	}

	calculatedSum := int64(0)
	for _, result := range calculatedResults {
		calculatedSum += result
	}

	fmt.Printf("Sum of calculated Cephalopod is %d\n", calculatedSum)
}

func initRegex() {
	whiteSpaceRegex = regexp.MustCompile(`(\ |\t)+`)
}

func parseInput(r io.Reader) ([]cephalopod.CephalopodRow, error) {
	sc := bufio.NewScanner(r)

	var inputTable [][]string
	for {
		if !sc.Scan() {
			if err := sc.Err(); err != nil {
				return nil, fmt.Errorf("erorr during parse input: %w", err)
			}
			break
		}

		line := sc.Text()
		line = strings.TrimSpace(line)
		line = whiteSpaceRegex.ReplaceAllString(line, ",")
		row := strings.Split(line, ",")

		if len(inputTable) != 0 && len(inputTable[0]) != len(row) {
			return nil, errors.New("rows should have all same width")
		}
		inputTable = append(inputTable, row)
	}

	cephalopodRows := make([]cephalopod.CephalopodRow, len(inputTable[0]))
	numRows := make([][]int64, len(inputTable[0]))

	for i, inputRow := range inputTable {
		isLastRow := i == len(inputTable)-1

		if !isLastRow {
			for j, nStr := range inputRow {
				n, err := strconv.ParseInt(nStr, 10, 64)
				if err != nil {
					return nil, fmt.Errorf("row should have valid number: %w", err)
				}

				numRows[j] = append(numRows[j], n)
			}
		} else {
			for j, opStr := range inputRow {
				op, err := operator.ParseOperator(opStr)
				if err != nil {
					return nil, fmt.Errorf("last row should have valid operator: %w", err)
				}

				cephalopodRows[j] = cephalopod.Builder().
					Nums(numRows[j]...).
					Operator(op).
					Build()
			}
		}
	}

	return cephalopodRows, nil
}
