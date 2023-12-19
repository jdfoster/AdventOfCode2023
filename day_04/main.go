package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var fn = "./day_04/input.txt"

func parseGame(s string) (int, []int, []int) {
	ss := strings.Split(s, ":")

	cs := strings.TrimSpace(strings.TrimPrefix(ss[0], "Card "))
	c, err := strconv.Atoi(cs)
	if err != nil {
		panic(err)
	}

	ll := strings.Split(ss[1], "|")
	r := make([][]int, len(ll))

	for i, l := range ll {
		gg := strings.Split(strings.TrimSpace(l), " ")
		r[i] = make([]int, 0, len(gg))

		for _, g := range gg {
			if len(g) > 0 {
				n, err := strconv.Atoi(g)
				if err != nil {
					panic(err)
				}

				r[i] = append(r[i], n)
			}
		}
	}

	return c, r[0], r[1]
}

func findMatched(aa, bb []int) []int {
	as := make(map[int]struct{}, len(aa))
	for _, a := range aa {
		as[a] = struct{}{}
	}

	rr := make([]int, 0, len(bb))
	for _, b := range bb {
		if _, ok := as[b]; ok {
			rr = append(rr, b)
		}
	}

	return rr
}

func calcScore(a []int) int {
	var v int
	l := len(a)

	if l > 0 {
		v = 1

	}

	for i := 0; i < l-1; i++ {
		v *= 2
	}

	return v
}

func sumCards(l string) int {
	_, a, b := parseGame(l)
	return calcScore(findMatched(a, b))
}

type counter struct {
	state map[int]int
}

func (c counter) Get(i int) int {
	if _, ok := c.state[i]; !ok {
		c.state[i] = 1
	}

	return c.state[i]
}

func (c counter) Add(l string) {
	i, a, b := parseGame(l)
	m := findMatched(a, b)

	if _, ok := c.state[i]; !ok {
		c.state[i] = 1
	}

	for j := 0; j < len(m); j++ {
		k := i + j + 1
		c.state[k] = c.Get(k) + c.Get(i)
	}
}

func (c counter) Sum() int {
	var r int
	for _, v := range c.state {
		r += v
	}

	return r
}

func main() {
	f, err := os.Open(fn)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	s := bufio.NewScanner(f)

	var a int
	b := &counter{state: make(map[int]int)}

	for s.Scan() {
		a += sumCards(s.Text())
		b.Add(s.Text())
	}

	fmt.Println("part one value: ", a)
	fmt.Println("part two value: ", b.Sum())
}
