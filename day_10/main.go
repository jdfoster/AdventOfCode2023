package main

import (
	"fmt"
	"os"
	"strings"
	"unicode/utf8"
)

var fn = "./day_10/input.txt"

type dir uint8

const (
	LEFT dir = iota
	UP
	RIGHT
	DOWN
)

type pipe struct {
	char   rune
	routes map[dir]dir
}

var (
	PIPE_HORIZONTAL pipe = pipe{
		char: '-',
		routes: map[dir]dir{
			LEFT:  LEFT,
			RIGHT: RIGHT,
		},
	}
	PIPE_VERTICAL pipe = pipe{
		char: '|',
		routes: map[dir]dir{
			UP:   UP,
			DOWN: DOWN,
		},
	}
	PIPE_UP_RIGHT pipe = pipe{
		char: 'F',
		routes: map[dir]dir{
			UP:   RIGHT,
			LEFT: DOWN,
		},
	}
	PIPE_UP_LEFT pipe = pipe{
		char: '7',
		routes: map[dir]dir{
			UP:    LEFT,
			RIGHT: DOWN,
		},
	}
	PIPE_DOWN_RIGHT pipe = pipe{
		char: 'L',
		routes: map[dir]dir{
			DOWN: RIGHT,
			LEFT: UP,
		},
	}
	PIPE_DOWN_LEFT pipe = pipe{
		char: 'J',
		routes: map[dir]dir{
			DOWN:  LEFT,
			RIGHT: UP,
		},
	}
	PIPE_START pipe = pipe{
		char: 'S',
		routes: map[dir]dir{
			LEFT:  LEFT,
			UP:    UP,
			RIGHT: RIGHT,
			DOWN:  DOWN,
		},
	}
	PIPE_NONE pipe = pipe{
		char:   '.',
		routes: map[dir]dir{},
	}
)

type position struct {
	i int
	d dir
}

type grid struct {
	cols  int
	rows  int
	start int
	pp    []pipe
}

func (g grid) line(i int, pp []position) string {
	u := g.cols * i
	cc := make(map[int]struct{}, g.cols)
	for _, p := range pp {
		if p.i >= u && p.i < u+g.cols {
			cc[p.i] = struct{}{}
		}
	}

	var b strings.Builder

	for i, p := range g.pp {
		if i < u || i >= u+g.cols {
			continue
		}

		_, found := cc[i]

		switch {
		case i == g.start:
			pi := g.findStartPipe()
			b.WriteRune(pi.char)
		case found:
			b.WriteRune(p.char)
		default:
			b.WriteRune('.')
		}
	}

	return b.String()
}

func (g grid) indexForDir(i int, d dir) (di int, ok bool) {
	ok = true
	switch d {
	case LEFT:
		di = i - 1
		ok = i%g.rows > 0
	case RIGHT:
		di = i + 1
		ok = i%g.rows < g.cols-1
	case UP:
		di = i - g.cols
	case DOWN:
		di = i + g.cols
	}

	ok = ok && di > -1 && di < g.cols*g.rows

	return
}

func (g grid) next(p position) (position, bool) {
	pi := g.pp[p.i]
	if di, ok := pi.routes[p.d]; ok {
		if j, ok := g.indexForDir(p.i, di); ok {
			pj := g.pp[j]
			_, ok = pj.routes[di]

			return position{j, di}, ok
		}
	}

	return p, false
}

func (g grid) findStartIndexes() ([]int, []dir) {
	jj := make([]int, 0, 4)
	dd := make([]dir, 0, 4)
	for d := range g.pp[g.start].routes {
		pi := position{g.start, d}
		if pj, ok := g.next(pi); ok {
			jj = append(jj, pj.i)
			dd = append(dd, pj.d)
		}
	}

	return jj, dd
}

func (g grid) findStartPosition() (int, position) {
	jj, dd := g.findStartIndexes()
	return jj[0], position{g.start, dd[1]}
}

func (g grid) findStartPipe() pipe {
	pp := []pipe{
		PIPE_HORIZONTAL,
		PIPE_VERTICAL,
		PIPE_UP_RIGHT,
		PIPE_UP_LEFT,
		PIPE_DOWN_RIGHT,
		PIPE_DOWN_LEFT,
	}
	_, dd := g.findStartIndexes()

	input, want := (dd[0] + 2%4), dd[1]
	for _, p := range pp {
		if got, ok := p.routes[input]; ok && got == want {
			return p
		}
	}

	panic("failed to find start")
}

func (g grid) walk() []position {
	var loop func(int, []position) []position

	loop = func(end int, walked []position) []position {
		if walked == nil || len(walked) < 1 {
			return walked
		}

		last := walked[len(walked)-1]

		if last.i == end {
			return walked
		}

		if next, ok := g.next(last); ok {
			return loop(end, append(walked, next))
		}

		panic("no route to end")
	}

	end, first := g.findStartPosition()
	rr := make([]position, 1, g.cols*g.rows)
	rr[0] = first

	return loop(end, rr)
}

func parse(s string) *grid {
	cols := strings.IndexByte(s, '\n')
	rows := (utf8.RuneCountInString(s) / cols) - 1

	g := &grid{
		cols: cols,
		rows: rows,
		pp:   make([]pipe, cols*rows),
	}

	var i int
	for _, r := range s {
		switch r {
		case '\n':
			continue
		case '-':
			g.pp[i] = PIPE_HORIZONTAL
		case '|':
			g.pp[i] = PIPE_VERTICAL
		case 'F':
			g.pp[i] = PIPE_UP_RIGHT
		case '7':
			g.pp[i] = PIPE_UP_LEFT
		case 'L':
			g.pp[i] = PIPE_DOWN_RIGHT
		case 'J':
			g.pp[i] = PIPE_DOWN_LEFT
		case 'S':
			g.pp[i] = PIPE_START
			g.start = i
		case '.':
			g.pp[i] = PIPE_NONE
		}

		i++
	}

	return g
}

func findEnclosed(s string) []int {
	rr := []rune(s)
	ee := make([]int, 0, len(rr))

	var (
		oi int
		o  bool
	)

	for i, r := range rr {

		switch r {
		case '|':
			o = !o
		case 'L':
			oi = i
		case 'F':
			oi = i
		case 'J':
			if rr[oi] == 'F' {
				o = !o
			}
		case '7':
			if rr[oi] == 'L' {
				o = !o
			}
		}

		if o && r == '.' {
			ee = append(ee, i)
		}
	}

	return ee
}

func main() {
	bb, err := os.ReadFile(fn)
	if err != nil {
		panic(err)
	}

	g := parse(string(bb))
	pp := g.walk()

	var b int
	for i := 0; i < g.rows; i++ {
		l := g.line(i, pp)
		li := []rune(l)
		for _, j := range findEnclosed(l) {
			li[j] = '*'
			b++
		}
		fmt.Println(string(li))
	}

	fmt.Println("part one value: ", len(pp)/2)
	fmt.Println("part two value: ", b)
}
