package main

import (
	"fmt"
	"os"
	"strings"
	"unicode/utf8"
)

var fn = "./day_10/input.txt"

type dir uint8

func (d dir) String() string {
	switch d {
	case LEFT:
		return "L"
	case UP:
		return "U"
	case RIGHT:
		return "R"
	case DOWN:
		return "D"
	}

	return "unknown"
}

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

func (g grid) printPositions(pp []position) string {
	dd := make(map[int]string, len(pp))

	for _, p := range pp {
		dd[p.i] = p.d.String()
	}

	var b strings.Builder

	for i, p := range g.pp {
		if i > 0 && i%g.cols == 0 {
			b.WriteRune('\n')
		}

		if d, ok := dd[i]; ok {
			b.WriteString(d)
			continue
		}

		switch p.char {
		case 'S':
			b.WriteRune('S')
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

func (g grid) walk(i int) []position {
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

	jj := make([]int, 0, 4)
	dd := make([]dir, 0, 4)
	for d := range g.pp[i].routes {
		pi := position{i, d}
		if pj, ok := g.next(pi); ok {
			jj = append(jj, pj.i)
			dd = append(dd, pj.d)
		}
	}

	if len(jj) != 2 || len(dd) != 2 {
		panic("must have start and finish")
	}

	rr := make([]position, 1, g.cols*g.rows)
	rr[0] = position{i, dd[1]}

	return loop(jj[0], rr)
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

func main() {
	bb, err := os.ReadFile(fn)
	if err != nil {
		panic(err)
	}

	g := parse(string(bb))
	w := g.walk(g.start)
	// fmt.Println(g.printPositions(w))
	fmt.Println("part one value: ", len(w)/2)
}
