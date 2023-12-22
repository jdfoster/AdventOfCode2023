package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

var fn = "./day_09/input.txt"

func parse(s string) []int {
	vv := strings.Split(s, " ")
	rr := make([]int, len(vv))

	for i, v := range vv {
		var err error
		rr[i], err = strconv.Atoi(v)
		if err != nil {
			panic(err)
		}
	}

	return rr
}

func diff(rr []int) []int {
	vv := make([]int, len(rr)-1)

	for i := 1; i < len(rr); i++ {
		vv[i-1] = rr[i] - rr[i-1]
	}

	return vv
}

func predict(rr []int) int {
	l := len(rr)
	p := rr[l-1]

	ri := rr
	for i := 1; i < l; i++ {
		ri = diff(ri)
		p += ri[len(ri)-1]
	}

	return p
}

func main() {
	f, err := os.Open(fn)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	s := bufio.NewScanner(f)

	var (
		a int
		b int
	)

	for s.Scan() {
		rr := parse(s.Text())
		a += predict(rr)
		slices.Reverse(rr)
		b += predict(rr)
	}

	fmt.Println("part one value: ", a)
	fmt.Println("part one value: ", b)
}
