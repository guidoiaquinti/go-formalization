package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"sort"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type recorder struct {
	trace []string
}

func (r *recorder) Len() int {
	r.trace = append(r.trace, "len 7")
	return 7
}
func (r *recorder) Less(i, j int) bool {
	result := i < j
	r.trace = append(r.trace, fmt.Sprintf("less %d %d %t", i, j, result))
	return result
}
func (r *recorder) Swap(i, j int) {
	r.trace = append(r.trace, fmt.Sprintf("swap %d %d", i, j))
}

type panickingLess struct {
	*recorder
	marker *int
}

func (r panickingLess) Less(i, j int) bool {
	r.recorder.Less(i, j)
	panic(r.marker)
}

func checkPanic() {
	marker := new(int)
	data := panickingLess{recorder: new(recorder), marker: marker}
	wrapped := sort.Reverse(data)
	defer func() {
		if p := recover(); p != marker || strings.Join(data.trace, "; ") != "less 2 1 false" {
			log.Fatalf("Less panic: recovered %v, trace %v", p, data.trace)
		}
		fmt.Println("Reverse: Less preserved the panic payload and prior trace.")
	}()
	wrapped.Less(1, 2)
}

func main() {
	checkPanic()
	check.CompareSamples("godiff-sort-reverse", [][2]int{{0, 0}, {0, 1}, {1, 0}, {math.MinInt, math.MaxInt}}, func(r *rand.Rand, _ int) [2]int {
		return [2]int{int(r.Uint64()), int(r.Uint64())}
	}, func(p [2]int) string {
		return fmt.Sprintf("%d %d", p[0], p[1])
	}, func(p [2]int) string {
		var data recorder
		reversed := sort.Reverse(&data)
		n := reversed.Len()
		less := reversed.Less(p[0], p[1])
		reversed.Swap(p[0], p[1])
		if n != 7 || less != (p[1] < p[0]) {
			return fmt.Sprintf("unexpected return values %d %t", n, less)
		}
		return strings.Join(data.trace, "; ")
	}, func(p [2]int) string {
		return fmt.Sprintf("len 7; less %d %d %t; swap %d %d", p[1], p[0], p[1] < p[0], p[0], p[1])
	})
}
