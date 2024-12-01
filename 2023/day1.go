package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
)

func Day1_half() {
    file, err := os.Open("input1.txt")
    if err != nil {
        fmt.Errorf("Error opening file: %v", err)
    }
    defer file.Close()

    zeroByte, nineByte := "0"[0], "9"[0]
    sum := 0
    for {
        var str string
        _, err := fmt.Fscanln(file, &str)
        if err == io.EOF {
            break
        } else {
            fmt.Errorf("Error reading file: %v", err)
        }

        var fstByte rune = -1
        var endByte rune
        for _, b := range str {
            if zeroByte <= byte(b) && byte(b) <= nineByte {
                if fstByte == -1 {
                    fstByte = b
                }
                endByte = b
            }
        }

        n, err := strconv.Atoi(string([]rune{fstByte, endByte}))
        if err != nil {
            fmt.Errorf("Error converting string to int: %v", err)
        }

        sum += n
    }

    fmt.Println(sum)
}

func Day1_full() {
    file, err := os.Open("input1.txt")
    if err != nil {
        fmt.Errorf("Error opening file: %v", err)
    }
    defer file.Close()

    zeroByte, nineByte := "0"[0], "9"[0]
    strmap := map[string]int{
        "zero": 0,
        "one": 1,
        "two": 2,
        "three": 3,
        "four": 4,
        "five": 5,
        "six": 6,
        "seven": 7,
        "eight": 8,
        "nine": 9,
    }
    sum := 0

    for {
        var str string
        _, err := fmt.Fscanln(file, &str)
        if err == io.EOF {
            break
        } else {
            fmt.Errorf("Error reading file: %v", err)
        }

        var fst int
        for i := 0; i < len(str); i++ {
            if zeroByte <= str[i] && str[i] <= nineByte {
                fst, _ = strconv.Atoi(str[i:i+1])
                break
            }
            if i + 3 <= len(str) {
                if intval, exists := strmap[str[i:i+3]]; exists {
                    fst = intval
                    break
                }
            }
            if i + 4 <= len(str) {
                if intval, exists := strmap[str[i:i+4]]; exists {
                    fst = intval
                    break
                }
            }
            if i + 5 <= len(str) {
                if intval, exists := strmap[str[i:i+5]]; exists {
                    fst = intval
                    break
                }
            }
        }

        var lst int
        for i := len(str)-1; i >= 0; i-- {
            if zeroByte <= str[i] && str[i] <= nineByte {
                lst, _ = strconv.Atoi(str[i:i+1])
                break
            }
            if i + 3 <= len(str) {
                if intval, exists := strmap[str[i:i+3]]; exists {
                    lst = intval
                    break
                }
            }
            if i + 4 <= len(str) {
                if intval, exists := strmap[str[i:i+4]]; exists {
                    lst = intval
                    break
                }
            }
            if i + 5 <= len(str) {
                if intval, exists := strmap[str[i:i+5]]; exists {
                    lst = intval
                    break
                }
            }
        }

        sum += fst*10 + lst
    }

    fmt.Println(sum)
}
