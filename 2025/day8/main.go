package main

import (
	"bufio"
	"container/heap"
	"errors"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"maps"
	"math"
	"os"
	"slices"
	"strings"

	"github.com/samber/lo"
)

var (
	filenameArg = flag.String("input", "", "input file path")
	connect     = flag.Int("connect", 0, "number of connections to make between boxes")
)

func main() {
	if err := parseArgs(); err != nil {
		slog.Error("invalid arguments", "error", err)
		os.Exit(1)
	}

	f, err := os.Open(*filenameArg)
	if err != nil {
		slog.Error("failed to open file", "error", err)
		os.Exit(1)
	}
	defer f.Close()

	boxes, err := parseBoxes(f)
	if err != nil {
		slog.Error("failed to parse boxes", "error", err)
		os.Exit(1)
	}

	// problem 1
	connectedBoxes := ConnectBoxes(slices.Clone(boxes), *connect)
	largestThreeMult := SizeMultipleOfLargestThreeCircuits(connectedBoxes)
	fmt.Printf("Size multiple of largest three circuits: %d\n", largestThreeMult)

	// problem 2
	lastConnectedPair := GetLastConnectPair(boxes)
	if lastConnectedPair == nil {
		slog.Warn("Input cannot be fully connected")
		os.Exit(1)
	}
	fmt.Printf("Multiply of X coordinates of last connected pair: %d\n",
		lastConnectedPair.Left.X*lastConnectedPair.Right.X)
}

func parseArgs() error {
	flag.Parse()

	if filenameArg == nil || *filenameArg == "" {
		return errors.New("filename should be given")
	}

	if connect == nil || *connect <= 0 {
		return errors.New("connect should be a positive integer")
	}

	return nil
}

func parseBoxes(r io.Reader) ([]Box, error) {
	sc := bufio.NewScanner(r)

	var boxes []Box

	for sc.Scan() {
		line := sc.Text()
		line = strings.TrimSpace(line)

		var x, y, z int
		if n, err := fmt.Sscanf(line, "%d,%d,%d", &x, &y, &z); err != nil {
			return nil, fmt.Errorf("failed to parse box from line %s: %w", line, err)
		} else if n < 3 {
			return nil, fmt.Errorf("invalid box format in line: %s", line)
		}

		boxes = append(boxes, Box{X: x, Y: y, Z: z})
	}

	if err := sc.Err(); err != nil {
		return nil, fmt.Errorf("error on scanning reader: %w", err)
	}

	return boxes, nil
}

func ConnectBoxes(boxes []Box, n int) [][]Box {
	circuits := initializeCircuits(boxes)
	hp := NewBoxPairDistanceHeap(boxes)
	heap.Init(hp)

	for range n {
		_ = connectBoxOnce(hp, &circuits)
	}

	return lo.Map(circuits, func(m map[Box]struct{}, _ int) []Box {
		return slices.Collect(maps.Keys(m))
	})
}

func GetLastConnectPair(boxes []Box) *BoxPair {
	totalBoxes := len(boxes)

	circuits := initializeCircuits(boxes)
	hp := NewBoxPairDistanceHeap(boxes)
	heap.Init(hp)

	for {
		boxPairPtr := connectBoxOnce(hp, &circuits)
		if boxPairPtr == nil {
			return nil
		}

		if len(circuits) == 1 && len(circuits[0]) == totalBoxes {
			return boxPairPtr
		}
	}
}

func initializeCircuits(boxes []Box) []map[Box]struct{} {
	return lo.Map(boxes, func(b Box, _ int) map[Box]struct{} {
		return map[Box]struct{}{
			b: {},
		}
	})
}

func connectBoxOnce(hp *BoxPairDistanceHeap, circuits *[]map[Box]struct{}) (connectedPair *BoxPair) {
	boxPairAny := heap.Pop(hp)
	if boxPairAny == nil {
		return nil
	}

	boxPair := boxPairAny.(BoxPair)

	var leftCircuitIdx, rightCircuitIdx *int

	for i, circuit := range *circuits {
		if _, ok := circuit[boxPair.Left]; ok {
			leftCircuitIdx = &i
		}
		if _, ok := circuit[boxPair.Right]; ok {
			rightCircuitIdx = &i
		}
	}

	switch {
	case leftCircuitIdx == nil || rightCircuitIdx == nil:
		panic("circuit initialization should have covered all boxes")
	case *leftCircuitIdx != *rightCircuitIdx:
		{
			maps.Copy((*circuits)[*leftCircuitIdx], (*circuits)[*rightCircuitIdx])
			*circuits = slices.Delete(*circuits, *rightCircuitIdx, *rightCircuitIdx+1)
		}
	default: // both in the same circuit
		{
		}
	}

	return &boxPair
}

func SizeMultipleOfLargestThreeCircuits(circuits [][]Box) int {
	hp := lo.ToPtr(CircuitMaxLenHeap(circuits))
	heap.Init(hp)

	mult := 1
	for range 3 {
		circuitAny := heap.Pop(hp)
		if circuitAny == nil {
			break
		}

		mult *= len(circuitAny.([]Box))
	}

	return mult
}

type Box struct {
	X, Y, Z int
}

type BoxPair struct {
	Left, Right Box
	distance    int
}

func MakeBoxPair(left Box, right Box) BoxPair {
	distance := int(math.Sqrt(
		math.Pow(math.Abs(float64(left.X-right.X)), 2) +
			math.Pow(math.Abs(float64(left.Y-right.Y)), 2) +
			math.Pow(math.Abs(float64(left.Z-right.Z)), 2),
	))

	return BoxPair{
		Left:     left,
		Right:    right,
		distance: distance,
	}
}

func (p BoxPair) Distance() int {
	return p.distance
}

type BoxPairDistanceHeap []BoxPair

func NewBoxPairDistanceHeap(boxes []Box) *BoxPairDistanceHeap {
	l := len(boxes)
	if l == 0 {
		return nil
	}

	var boxPairs []BoxPair
	for i := range l {
		for j := i + 1; j < l; j++ {
			boxPair := MakeBoxPair(boxes[i], boxes[j])
			boxPairs = append(boxPairs, boxPair)
		}
	}

	hp := BoxPairDistanceHeap(boxPairs)
	heap.Init(&hp)
	return &hp
}

func (b *BoxPairDistanceHeap) Len() int {
	return len(*b)
}

func (b *BoxPairDistanceHeap) Less(i int, j int) bool {
	lp := (*b)[i]
	rp := (*b)[j]

	return lp.Distance() < rp.Distance()
}

func (b *BoxPairDistanceHeap) Pop() any {
	if b.Len() == 0 {
		return nil
	}

	elem := (*b)[b.Len()-1]
	*b = (*b)[:b.Len()-1]
	return elem
}

func (b *BoxPairDistanceHeap) Push(x any) {
	(*b) = append(*b, x.(BoxPair))
}

func (b *BoxPairDistanceHeap) Swap(i int, j int) {
	(*b)[i], (*b)[j] = (*b)[j], (*b)[i]
}

type CircuitMaxLenHeap [][]Box

func (c *CircuitMaxLenHeap) Len() int {
	return len(*c)
}

func (c *CircuitMaxLenHeap) Less(i int, j int) bool {
	iLen := len((*c)[i])
	jLen := len((*c)[j])

	return jLen < iLen // since we want max-heap
}

func (c *CircuitMaxLenHeap) Pop() any {
	if c.Len() == 0 {
		return nil
	}

	elem := (*c)[c.Len()-1]
	(*c) = (*c)[:c.Len()-1]
	return elem
}

func (c *CircuitMaxLenHeap) Push(x any) {
	*c = append(*c, x.([]Box))
}

func (c *CircuitMaxLenHeap) Swap(i int, j int) {
	(*c)[i], (*c)[j] = (*c)[j], (*c)[i]
}
