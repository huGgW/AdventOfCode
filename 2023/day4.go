package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"slices"
	"strconv"
	"strings"
	"math"
)

func Day4Half() {
    f, err := os.Open("input4.txt")
    if err != nil {
        panic(err)
    }
    defer f.Close()
    reader := bufio.NewReader(f)

    totalPoint := 0
    for {
        line, err := reader.ReadString('\n')
        if err == io.EOF {
            break
        } else if err != nil {
            panic(err)
        }

        gameSlice, playerSlice := getCards(line)
        point := calcPoint(gameSlice, playerSlice)
        totalPoint += point
    }

    fmt.Println(totalPoint)
}


func Day4Full() {
    f, err := os.Open("input4.txt")
    if err != nil {
        panic(err)
    }
    defer f.Close()
    reader := bufio.NewReader(f)

    cards := []int{0}
    lineIdx := 0
    for {
        line, err := reader.ReadString('\n')
        if err == io.EOF {
            break
        } else if err != nil {
            panic(err)
        }

        if (len(cards) <= lineIdx) {
            cards = append(cards, 1)
        } else {
            cards[lineIdx]++
        }
        games, players := getCards(line)

        duplicateCards(lineIdx, games, players, &cards)

        lineIdx++
    }

    sum := 0
    for _, cnt := range cards {
        sum += cnt
    }
    fmt.Println(sum)
}

func getCards(line string) ([]int, []int) {
    spl := strings.Index(line, ":")
    cards := strings.Split(line[spl+1:], " | ")

    gamesStr := strings.Split(cards[0], " ")
    playersStr := strings.Split(cards[1], " ")

    games := []int{}
    for _, gs := range gamesStr {
        gsTrim := strings.TrimSpace(gs)
        if gsTrim == "" {
            continue
        }
        game, err := strconv.Atoi(gsTrim)
        if err != nil {
            panic(err)
        }
        games = append(games, game)
    }

    players := []int{}
    for _, ps := range playersStr {
        psTrim := strings.TrimSpace(ps)
        if psTrim == "" {
            continue
        }
        player, err := strconv.Atoi(psTrim)
        if err != nil {
            panic(err)
        }
        players = append(players, player)
    }

    return games, players
}

func cntSameCards(games, players []int) int {
    gamesSorted := slices.Clone(games)
    slices.Sort(gamesSorted)

    cnt := 0
    for _, player := range players {
        _, found := slices.BinarySearch(gamesSorted, player)
        if found {
            cnt++
        }
    }

    return cnt
}

func calcPoint(games, players []int) int {
    cnt := cntSameCards(games, players)

    if cnt == 0 {
        return 0
    } else {
        return int(math.Pow(2.0, float64(cnt-1)))
    }
}

func duplicateCards(line int, games, players []int, cards *([]int)) {
    cnt := cntSameCards(games, players)

    for idx := line+1; idx <= line + cnt; idx++ {
        if (idx >= len(*cards)) {
            *cards = append(*cards, (*cards)[line])
        } else {
            (*cards)[idx] += (*cards)[line]
        }
    }
}
