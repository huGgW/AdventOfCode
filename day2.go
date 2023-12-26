package main

import (
	"fmt"
    "bufio"
	"os"
	"strconv"
	"strings"
)


type gameStat struct {
    red int
    blue int
    green int
}
func gameStatFromString(str string) gameStat {
    newGameState := gameStat{red: 0, blue: 0, green: 0}

    entrySlice := strings.Split(str, ",")
    for _, entry := range entrySlice {

        splitEntry := strings.Split(strings.TrimSpace(entry), " ")
        num, _ := strconv.Atoi(splitEntry[0])
        switch splitEntry[1] {
        case "blue": 
            newGameState.blue = num
        case "red":
            newGameState.red = num
        case "green":
            newGameState.green = num
        }
    }

    return newGameState
}

type totalCnt gameStat
func (tot totalCnt) isAvailable(gm gameStat) bool {
    return tot.red >= gm.red && tot.blue >= gm.blue && tot.green >= gm.green
}
func (totPtr *totalCnt) makePossible(gm gameStat) {
    if totPtr.red < gm.red {
        totPtr.red = gm.red
    }
    if totPtr.blue < gm.blue {
        totPtr.blue = gm.blue
    }
    if totPtr.green < gm.green {
        totPtr.green = gm.green
    }
}
func (tot totalCnt) power() int {
    return tot.red * tot.blue * tot.green
}


func analyzeGameInput(str string) (int, []gameStat) {
    spl := strings.Split(str, ":")

    id, _ := strconv.Atoi(strings.Split(spl[0], " ")[1])

    gameStatSlice := []gameStat{}
    statsStrSlice := strings.Split(spl[1], ";")
    for _, statStr := range statsStrSlice {
        gameStatSlice = append(gameStatSlice, gameStatFromString(statStr))
    }

    return id, gameStatSlice
}

func Day2Half() {
    f, err := os.Open("input2.txt")
    if err != nil {
        panic(fmt.Errorf("Error opening file: %v", err))
    }
    defer f.Close()
    scanner := bufio.NewScanner(f)

    total := totalCnt{red: 12, green: 13, blue: 14}
    sum := 0

    for {
        if !scanner.Scan() {
            err := scanner.Err()
            if err == nil { // EOF
                break
            } else {
                panic(fmt.Errorf("Error reading file: %v", err))
            }
        }
        str := scanner.Text()

        id, gameStatSlice := analyzeGameInput(str)

        gameAvailable := true
        for _, game := range gameStatSlice {
            if !total.isAvailable(game) {
                gameAvailable = false
                break
            }
        }
        if gameAvailable {
            sum += id
        }
    }

    fmt.Println(sum)
}

func Day2Full() {
    f, err := os.Open("input2.txt")
    if err != nil {
        panic(fmt.Errorf("Error opening file: %v", err))
    }
    defer f.Close()
    scanner := bufio.NewScanner(f)

    sumOfPower := 0
    for {
        if !scanner.Scan() {
            err := scanner.Err()
            if err == nil { // EOF
                break
            } else {
                panic(fmt.Errorf("Error reading file: %v", err))
            }
        }
        str := scanner.Text()

        availableTot := totalCnt{0, 0, 0}

        _, gameStatSlice := analyzeGameInput(str)
        for _, game := range gameStatSlice {
            availableTot.makePossible(game)
        }

        sumOfPower += availableTot.power()
    }

    fmt.Println(sumOfPower)
}

