package main

import (
    "fmt"
    "os"
    "strconv"
)

type numPos struct {
    row int
    start int
    end int
}

type symbolPos struct {
    row int
    col int
}

func (sym *symbolPos)isNeighborTo(num *numPos) bool {
    if !(sym.row - 1 <= num.row && num.row <= sym.row + 1) {
        return false
    }

    if !(sym.col - 1 <= num.end && num.start <= sym.col + 1) {
        return false
    }

    return true
}

func (sym *symbolPos)analyzeGear(schem [][]rune, nums []numPos) (int, bool) {
    if schem[sym.row][sym.col] != '*' {
        return 0, false
    }

    neighbors := []numPos{}
    for _, num := range nums {
        if sym.isNeighborTo(&num) {
            neighbors = append(neighbors, num)
        }

        if len(neighbors) > 2 {
            return 0, false
        }
    }

    if len(neighbors) == 2 {
        x, err := strconv.Atoi(string(schem[neighbors[0].row][neighbors[0].start:neighbors[0].end + 1]))
        if err != nil {
            panic(err)
        }
        y, err := strconv.Atoi(string(schem[neighbors[1].row][neighbors[1].start:neighbors[1].end + 1]))
        if err != nil {
            panic(err)
        }

        return x * y, true
    } else {
        return 0, false
    }
}


func Day3Half() {
    schem, nums, symbols := getFromFile("input3.txt")

    // debugPrint(schem, nums, symbols)
    // fmt.Println()

    partNums := NewSet[numPos]()

    for _, num := range nums.ToSlice() {
        for _, symbol := range symbols.ToSlice() {
            if symbol.isNeighborTo(&num) {
                partNums.Add(num)
            }
        }
    }

    sum := 0
    for _, numP := range partNums.ToSlice() {
        n, err := strconv.Atoi(string(schem[numP.row][numP.start:numP.end + 1]))
        if err != nil {
            panic(err)
        }
        // fmt.Printf("......%d\n", n)
        sum += n
    }

    fmt.Println(sum)
}

func Day3Full() {
    schem, nums, symbols := getFromFile("input3.txt")

    gearNums := []int{}
    numSlice := nums.ToSlice()

    for _, symbol := range symbols.ToSlice() {
        gearNum, isGear := symbol.analyzeGear(schem, numSlice)
        if isGear {
            gearNums = append(gearNums, gearNum)
        }
    }

    result := 0
    for _, gearNum := range gearNums {
        result += gearNum
    }

    fmt.Println(result)
}


func getFromFile(fname string) ([][]rune, *Set[numPos], *Set[symbolPos]) {
    f, err := os.Open(fname)
    if err != nil {
        panic(err)
    }

    defer f.Close()

    var str string


    n, err := fmt.Fscanln(f, &str)
    if err != nil {
        panic(err)
    }

    schem := [][]rune{}
    nums := NewSet[numPos]()
    symbols := NewSet[symbolPos]()

    row := 0
    for ; n > 0; n, err = fmt.Fscanln(f, &str) {
        if err != nil {
            panic(err)
        }

        lines := []rune(str)
        schem = append(schem, lines)

        start, end := -1, -1
        for i, r := range lines {
            switch {
            case r == '.': {
                if end != -1 {
                    nums.Add(numPos{row, start, end})
                    start, end = -1, -1
                }
            }
            case '0' <= r && r <= '9': {
                if start == -1 {
                    start = i
                }
                end = i
            }
            default: {
                if end != -1 {
                    nums.Add(numPos{row, start, end})
                    start, end = -1, -1
                }
                symbols.Add(symbolPos{row, i})
            }
            }
        }
        if end != -1 {
            nums.Add(numPos{row, start, end})
        }

        row++
    }

    return schem, nums, symbols
}

func debugPrint(schemes [][]rune, nums *Set[numPos], symbols *Set[symbolPos]) {

    fmt.Println("Numbers:")
    for _, num := range nums.ToSlice() {
        fmt.Printf("\t%dth row: %d ~ %d\n", num.row, num.start, num.end)
    }

    fmt.Println("Symbols:")
    for _, symbol := range symbols.ToSlice() {
        fmt.Printf("\trow: %d, column: %d\n", symbol.row, symbol.col)
    }

    fmt.Println("Sheme:")
    for _, scheme := range schemes {
        fmt.Printf("\t%s\n", string(scheme))
    }
}
