package grid

import (
	"day4/pair"
	"fmt"
)

type Grid struct {
	grid [][]Item
}

func (g *Grid) RemovePapers(positions []pair.Pair) error {
	if len(g.grid) < 2 || len(g.grid[0]) < 2 {
		panic("invalid grid")
	}

	h := len(g.grid)
	w := len(g.grid[0])

	for _, pos := range positions {
		if pos.I <= 0 || pos.I >= h-1 ||
			pos.J <= 0 || pos.J >= w-1 {
			return fmt.Errorf("cannot remove paper at edge or invalid position: (%d, %d)", pos.I, pos.J)
		}

		if g.grid[pos.I][pos.J] != Paper {
			return fmt.Errorf("given position is not a paper: (%d, %d)", pos.I, pos.J)
		}

		g.grid[pos.I][pos.J] = Empty
	}

	return nil
}

func (g *Grid) ForkLiftAccessablePosition() []pair.Pair {
	if len(g.grid) < 2 || len(g.grid[0]) < 2 {
		panic("invalid grid")
	}

	var accessable []pair.Pair
	h := len(g.grid)
	w := len(g.grid[0])

	for i := 1; i < h-1; i++ {
		for j := 1; j < w-1; j++ {
			if g.grid[i][j] != Paper {
				continue
			}

			if g.countAdjacentPapers(i, j) < 4 {
				accessable = append(accessable, pair.Pair{i, j})
			}
		}
	}

	return accessable
}

func (g *Grid) countAdjacentPapers(i, j int) uint8 {
	if i == 0 || j == 0 || i == len(g.grid)-1 || j == len(g.grid[0])-1 {
		panic("edge position should not be checked")
	}

	var count uint8
	for di := -1; di <= 1; di++ {
		for dj := -1; dj <= 1; dj++ {
			if di == 0 && dj == 0 {
				continue
			}

			if g.grid[i+di][j+dj] == Paper {
				count++
			}
		}
	}

	return count
}

type Item uint8

const (
	Invalid Item = iota
	Empty
	Paper
)

func ParseItem(r rune) Item {
	switch r {
	case '@':
		return Paper
	case '.':
		return Empty
	default:
		return Invalid
	}
}
