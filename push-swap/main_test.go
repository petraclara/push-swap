package main

import (
	"fmt"
	"math/rand"
	"testing"
	"push-swap/stack"
)

func TestMoveCounts(t *testing.T) {
	sizes := []int{100, 500}
	limits := []int{700, 5500}

	for i, size := range sizes {
		perm := rand.Perm(size)
		s := stack.Stack(perm)
		a := stack.Indexify(s)
		b := stack.Stack{}

		// We need to capture the output, but since we just want to count, 
		// we can calculate moves by running the algorithm without printing.
		moves := countMoves(a, b)
		fmt.Printf("Size %d: %d moves (Limit %d)\n", size, moves, limits[i])
		if moves >= limits[i] {
			t.Errorf("Size %d used %d moves, which is >= %d", size, moves, limits[i])
		}
	}
}

func countMoves(a, b stack.Stack) int {
	moves := 0
	size := len(a)
	chunkSize := 15
	if size > 100 {
		chunkSize = 35
	}

	for i := 0; len(a) > 0; {
		if a[0] <= i {
			stack.Pb(&a, &b)
			moves++
			stack.Rb(&b)
			moves++
			i++
		} else if a[0] <= i+chunkSize {
			stack.Pb(&a, &b)
			moves++
			i++
		} else {
			stack.Ra(&a)
			moves++
		}
	}

	for len(b) > 0 {
		maxIdx := 0
		for i, v := range b {
			if v > b[maxIdx] {
				maxIdx = i
			}
		}

		if maxIdx <= len(b)/2 {
			for i := 0; i < maxIdx; i++ {
				stack.Rb(&b)
				moves++
			}
		} else {
			for i := 0; i < len(b)-maxIdx; i++ {
				stack.Rrb(&b)
				moves++
			}
		}
		stack.Pa(&a, &b)
		moves++
	}
	return moves
}
