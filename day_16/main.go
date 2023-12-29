package main

import (
	"fmt"
	"os"
	"strings"
	"unicode/utf8"
)

var fn = "./day_16/input.txt"

type dir uint8

const (
	LEFT dir = iota
	UP
	RIGHT
	DOWN
)

type grid struct {
	cols int
	rows int
	rr   []rune
}

func (g grid) String() string {
	var b strings.Builder

	for i, r := range g.rr {
		if i > 0 && i%g.cols == 0 {
			b.WriteRune('\n')
		}
		b.WriteRune(r)
	}

	return b.String()
}

func (g grid) countRune(r rune) int {
	var c int

	for _, ri := range g.rr {
		if ri == r {
			c++
		}
	}

	return c
}

func (g *grid) fillRune(r rune) {
	l := g.cols * g.rows

	for i := 0; i < l; i++ {
		g.rr[i] = r
	}
}

func (g *grid) fillString(s string) {
	var i int
	for _, r := range s {
		if r != '\n' {
			g.rr[i] = r
			i++
		}
	}
}

func newGrid(c, r int) *grid {
	l := c * r

	return &grid{
		cols: c,
		rows: r,
		rr:   make([]rune, l),
	}
}

func newGridFromString(s string) *grid {
	c := strings.IndexRune(s, '\n')
	r := (utf8.RuneCountInString(s) / c) - 1
	g := newGrid(c, r)
	g.fillString(s)

	return g
}

func newGridEmpty(c, r int, v rune) *grid {
	g := newGrid(c, r)
	g.fillRune(v)

	return g
}

type position struct {
	i int
	d dir
}

func (p position) dir(g *grid) []dir {
	r := make([]dir, 0, 2)
	v := g.rr[p.i]

	switch {
	case v == '\\' && p.d == LEFT:
		r = append(r, UP)
	case v == '\\' && p.d == UP:
		r = append(r, LEFT)
	case v == '\\' && p.d == RIGHT:
		r = append(r, DOWN)
	case v == '\\' && p.d == DOWN:
		r = append(r, RIGHT)
	case v == '/' && p.d == LEFT:
		r = append(r, DOWN)
	case v == '/' && p.d == UP:
		r = append(r, RIGHT)
	case v == '/' && p.d == RIGHT:
		r = append(r, UP)
	case v == '/' && p.d == DOWN:
		r = append(r, LEFT)
	case v == '|' && (p.d == LEFT || p.d == RIGHT):
		r = append(r, UP, DOWN)
	case v == '-' && (p.d == UP || p.d == DOWN):
		r = append(r, LEFT, RIGHT)
	default:
		r = append(r, p.d)
	}

	return r
}

func (p position) move(g *grid, d dir) (di int, ok bool) {
	ok = true
	switch d {
	case LEFT:
		di = p.i - 1
		ok = p.i%g.rows > 0
	case RIGHT:
		di = p.i + 1
		ok = p.i%g.rows < g.cols-1
	case UP:
		di = p.i - g.cols
	case DOWN:
		di = p.i + g.cols
	}

	ok = ok && di > -1 && di < g.cols*g.rows

	return di, ok
}

func walk(g *grid, p position) []position {
	state := make(map[position]struct{}, g.cols*g.rows)
	state[p] = struct{}{}

	var wrapper func(position)

	wrapper = func(q position) {
		for _, d := range q.dir(g) {
			if i, ok := q.move(g, d); ok {
				qi := position{i: i, d: d}
				_, known := state[qi]
				if !known {
					state[qi] = struct{}{}
					wrapper(qi)
				}
			}
		}
	}

	wrapper(p)

	pp := make([]position, 0, len(state))

	for q := range state {
		pp = append(pp, q)
	}

	return pp
}

func writeRoute(g *grid, pp []position) *grid {
	r := newGridEmpty(g.cols, g.rows, '.')

	for _, p := range pp {
		r.rr[p.i] = '#'
	}

	return r
}

func findEdges(c, r int) []position {
	l := r * c
	rr := make([]position, 0, l)

	for i := 0; i < c; i++ {
		rr = append(rr, position{i: i, d: DOWN})
	}

	for i := l - r; i < l; i++ {
		rr = append(rr, position{i: i, d: UP})
	}

	for i := 0; i < l; i++ {
		if i%r == 0 {
			rr = append(rr, position{i: i, d: RIGHT})
		}
		if i%r == r-1 {
			rr = append(rr, position{i: i, d: LEFT})
		}
	}

	return rr
}

func main() {
	bb, err := os.ReadFile(fn)
	if err != nil {
		panic(err)
	}

	g := newGridFromString(string(bb))

	pp := walk(g, position{i: 0, d: RIGHT})
	a := writeRoute(g, pp)
	// fmt.Println(a)
	// fmt.Println("")
	fmt.Println("part one value: ", a.countRune('#'))
	// fmt.Println("")

	dd := findEdges(g.cols, g.rows)
	gg := make([]*grid, len(dd))
	var mx, mi int
	for i, d := range dd {
		w := walk(g, d)
		gg[i] = writeRoute(g, w)
		if c := gg[i].countRune('#'); c > mx {
			mx = c
			mi = i
		}
	}

	b := gg[mi]
	// fmt.Println(b)
	// fmt.Println("")
	fmt.Println("part two value: ", b.countRune('#'))
}
