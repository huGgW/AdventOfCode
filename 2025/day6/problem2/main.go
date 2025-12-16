package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strconv"
	"unicode"

	"day6/cephalopod"
	"day6/operator"
	"day6/util"
)

func main() {
	filePath := util.InitFlags()
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("failed to open input file: %v", err)
	}
	defer f.Close()

	rs, err := readInput(f)
	if err != nil {
		log.Fatalf("failed to parse file: %v", err)
	}

	if err := validateInput(rs); err != nil {
		log.Fatalf("input is not valid: %v", err)
	}

	cephalopodRows := parseCephalopodRows(rs)

	calculatedResults := make([]int64, 0, len(cephalopodRows))
	for _, c := range cephalopodRows {
		x := c.Calculate()
		calculatedResults = append(calculatedResults, x)
	}

	var calculatedSum int64
	for _, x := range calculatedResults {
		calculatedSum += x
	}

	fmt.Printf("Sum of vertically calculated Cephalopod is %d\n", calculatedSum)
}

func readInput(r io.Reader) ([][]rune, error) {
	sc := bufio.NewScanner(r)

	var runes [][]rune

	for sc.Scan() {
		line := sc.Text()
		runes = append(runes, []rune(line))
	}

	if err := sc.Err(); err != nil {
		return nil, fmt.Errorf("failed to read input: %w", err)
	}

	return runes, nil
}

func validateInput(rs [][]rune) error {
	// should have at least two row
	if len(rs) < 2 {
		return fmt.Errorf("should have at least 2 row")
	}

	// last row should have operators
	for _, oRune := range rs[len(rs)-1] {
		if unicode.IsSpace(oRune) {
			continue
		}

		if _, err := operator.ParseOperator(string(oRune)); err == nil {
			continue
		}

		return fmt.Errorf("last row should have only operators")
	}

	// all operator should have at least one number in column
	for j, oRune := range rs[len(rs)-1] {
		if unicode.IsSpace(oRune) {
			continue
		}

		var hasDigit bool
		for _, row := range rs {
			dRune := row[j]
			if unicode.IsDigit(dRune) {
				hasDigit = true
				break
			}
		}

		if !hasDigit {
			return fmt.Errorf("column of operator should has at least one digit")
		}
	}

	return nil
}

func parseCephalopodRows(rs [][]rune) []cephalopod.CephalopodRow {
	h := len(rs)

	var cephalopodRows []cephalopod.CephalopodRow
	latestHandledCol := len(rs[h-1])

	for latestHandledCol > 0 {
		cephalopodRow, col := parseSingleCaphalopodRow(rs, latestHandledCol-1)
		cephalopodRows = append(cephalopodRows, cephalopodRow)
		latestHandledCol = col
	}

	return cephalopodRows
}

func parseSingleCaphalopodRow(rs [][]rune, from int) (cephalopod.CephalopodRow, int) {
	h := len(rs)
	operatorRow := rs[h-1]

	col := from
	builder := cephalopod.Builder()

	for {

		var digits []int
		for row := 0; row < h-1; row++ {
			r := rs[row][col]

			switch {
			case unicode.IsSpace(r):
				continue
			case unicode.IsDigit(r):
				d, err := strconv.Atoi(string(r))
				if err != nil {
					panic(err)
				}

				digits = append(digits, d)
			default:
				panic("sholud not be reached")
			}
		}

		if len(digits) > 0 {
			builder.Nums(digitsToNumber(digits))
		}

		if unicode.IsSpace(operatorRow[col]) {
			col--
			continue
		}

		op, err := operator.ParseOperator(string(operatorRow[col]))
		if err != nil {
			panic(err)
		}

		builder.Operator(op)
		break
	}

	return builder.Build(), col
}

func digitsToNumber(digits []int) int64 {
	x := int64(0)

	for i, d := range digits {
		x += int64(math.Pow10(len(digits)-i-1)) * int64(d)
	}

	return x
}
