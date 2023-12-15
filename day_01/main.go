package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

var fn = "./day_01/input.txt"

type position struct {
	i int
	r rune
}

func findNumericalValues(s string) (pp []position) {
	for i, r := range s {
		if unicode.IsNumber(r) {
			pp = append(pp, position{i: i, r: r})
		}
	}

	return pp
}

var re = regexp.MustCompile("one|two|three|four|five|six|seven|eight|nine")

func findWordValue(s string) (pp []position) {
	b := []byte(strings.ToLower(s))

	i := 0

	for {
		v := re.FindIndex(b[i:])
		if v == nil {
			break
		}

		ai, bi := v[0]+i, v[1]+i

		p := position{i: ai}
		i = p.i + 1

		switch w := string(b[ai:bi]); w {
		case "one":
			p.r = '1'
		case "two":
			p.r = '2'
		case "three":
			p.r = '3'
		case "four":
			p.r = '4'
		case "five":
			p.r = '5'
		case "six":
			p.r = '6'
		case "seven":
			p.r = '7'
		case "eight":
			p.r = '8'
		case "nine":
			p.r = '9'
		default:
			panic(fmt.Sprintf("failed to match %q", w))
		}

		pp = append(pp, p)
	}

	return pp
}

func runesToNumber(rr ...rune) int {
	v, err := strconv.Atoi(string(rr))
	if err != nil {
		panic(err)
	}

	return v
}

func numericalValue(s string) int {
	nn := findNumericalValues(s)
	return runesToNumber(nn[0].r, nn[len(nn)-1].r)
}

func characterValue(s string) int {
	nn := findNumericalValues(s)
	a, b := nn[0], nn[len(nn)-1]

	ww := findWordValue(s)
	if l := len(ww); l > 0 {
		if w := ww[0]; a.i > w.i {
			a = w
		}

		if w := ww[l-1]; b.i < w.i {
			b = w
		}
	}

	return runesToNumber(a.r, b.r)
}

func main() {
	f, err := os.Open(fn)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	s := bufio.NewScanner(f)

	var a, b int

	for s.Scan() {
		a += numericalValue(s.Text())
		b += characterValue(s.Text())
	}

	fmt.Println("part one value: ", a)
	fmt.Println("part two value: ", b)
}
