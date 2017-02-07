package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"unicode"
)

const (
	INPUT_FILENAME = "input.txt"
)

func main() {
	f, err := os.Open(INPUT_FILENAME)
	if err != nil {
		log.Fatalln("Failed to open input file:", err)
	}
	defer f.Close()

	r := bufio.NewReader(f)
	// ignore errors because of task convention
	line, _ := r.ReadString('\n')
	flprt := strings.Split(line, " ")
	m, _ := strconv.Atoi(flprt[0])
	// n, _ := strconv.Atoi(flprt[1])

	in := make([][]rune, m)
	for i, _ := range in {
		l, _ := r.ReadString('\n')
		in[i] = []rune(l)
	}

	key := "OneTwoTrip"

	type Match struct {
		r rune
		x int
		y int
	}
	out := make([]Match, 0, len(key))

	for _, s := range key {
		s = unicode.ToLower(s)
		found := false
		for i, a := range in {
			for j, r := range a {
				if unicode.ToLower(r) == s {
					out = append(out, Match{r, i, j})
					in[i][j] = 0
					found = true
					break
				}
			}
			if found {
				break
			}
		}
		if !found {
			fmt.Println("Impossible")
			os.Exit(1)
		}
	}

	for _, m := range out {
		fmt.Printf("%s - (%d, %d);\n", string(m.r), m.x, m.y)
	}
}
