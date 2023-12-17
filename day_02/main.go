package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var (
	fn     = "./day_02/input.txt"
	inital = map[string]int{
		"red":   12,
		"green": 13,
		"blue":  14,
	}
)

func parseGame(s string) (int, []map[string]int) {
	gg := strings.Split(strings.TrimPrefix(s, "Game "), ":")
	g, err := strconv.Atoi(gg[0])
	if err != nil {
		panic(err)
	}

	pp := strings.Split(gg[1], ";")
	aa := make([]map[string]int, len(pp))
	for i, p := range pp {
		rr := strings.Split(strings.TrimSpace(p), ",")
		aa[i] = make(map[string]int, len(rr))

		for _, r := range rr {
			v := strings.Split(strings.TrimSpace(r), " ")
			c, err := strconv.Atoi(v[0])
			if err != nil {
				panic(err)
			}

			aa[i][strings.ToLower(v[1])] = c
		}
	}

	return g, aa
}

func main() {
	f, err := os.Open(fn)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	s := bufio.NewScanner(f)

	var a int

LOOP:
	for s.Scan() {
		g, rr := parseGame(s.Text())

		for _, r := range rr {
			for c, v := range r {
				w, ok := inital[c]
				if !ok {
					continue LOOP
				}

				if v > w {
					continue LOOP
				}
			}
		}

		a += g
	}

	fmt.Println("part one value: ", a)
}
