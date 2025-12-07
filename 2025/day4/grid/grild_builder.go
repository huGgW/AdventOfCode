package grid

import "errors"

type GridBuilder struct {
	rawGrid [][]rune
}

func (gb *GridBuilder) Append(runes []rune) error {
	if len(gb.rawGrid) > 0 && len(gb.rawGrid[0]) != len(runes) {
		return errors.New("inconsistent row length")
	}

	gb.rawGrid = append(gb.rawGrid, runes)
	return nil
}

func (gb *GridBuilder) Build() (Grid, error) {
	rawHeight := len(gb.rawGrid)
	if rawHeight == 0 {
		return Grid{}, nil
	}

	rawWidth := len(gb.rawGrid[0])

	grid := make([][]Item, 0, rawHeight+2)

	grid = append(grid, gb.invalidRow(rawWidth+2))

	for _, rawRow := range gb.rawGrid {
		row := make([]Item, 0, rawWidth+2)
		row = append(row, Invalid)

		for _, r := range rawRow {
			item := ParseItem(r)
			if item == Invalid {
				return Grid{}, errors.New("invalid item in input")
			}
			row = append(row, item)
		}

		row = append(row, Invalid)

		grid = append(grid, row)
	}

	grid = append(grid, gb.invalidRow(rawWidth+2))

	return Grid{grid}, nil
}

func (gb *GridBuilder) invalidRow(width int) []Item {
	r := make([]Item, 0, width)

	for range width {
		r = append(r, Invalid)
	}

	return r
}
