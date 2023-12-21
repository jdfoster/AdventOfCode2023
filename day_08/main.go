package main

import (
	"fmt"
	"os"
	"strings"
	"unicode/utf8"
)

var fn = "./day_08/input.txt"

func parse(s string) (string, [][]string) {
	vv := strings.Split(s, "\n\n")
	ll := strings.Split(vv[1], "\n")
	ee := make([][]string, 0, len(ll))

	for _, l := range ll {
		if l == "" {
			continue
		}

		e := make([]string, 3)
		e[0] = l[:3]
		e[1] = l[7:10]
		e[2] = l[12:15]

		ee = append(ee, e)
	}

	return vv[0], ee
}

type node struct {
	left  string
	right string
}

type chart struct {
	directions string
	nodes      map[string]node
}

func (c chart) move(p string, r rune) string {
	n, ok := c.nodes[p]
	if !ok {
		panic("node not found: " + p)
	}

	switch r {
	case 'L':
		return n.left
	case 'R':
		return n.right
	}

	panic("unknown direction: " + string(r))
}

func (c chart) next(i int) rune {
	j := i % utf8.RuneCount([]byte(c.directions))
	return []rune(c.directions)[j]
}

func (c chart) run(start, end string) int {
	p := start
	var i int

	for p != end {
		p = c.move(p, c.next(i))
		i++
	}

	return i
}

func newChart(d string, ee [][]string) *chart {
	c := &chart{
		directions: d,
		nodes:      make(map[string]node, len(ee)),
	}

	for _, e := range ee {
		c.nodes[e[0]] = node{
			left:  e[1],
			right: e[2],
		}
	}

	return c
}

func main() {
	bb, err := os.ReadFile(fn)
	if err != nil {
		panic(err)
	}

	d, ee := parse(string(bb))
	c := newChart(d, ee)

	fmt.Println("part one value: ", c.run("AAA", "ZZZ"))
}
