package main

import (
	"fmt"
	"os"

	"push-swap/stack"
)

func main() {
	// No arguments → print nothing
	if len(os.Args) < 2 {
		return
	}

	// Parse input
	a, err := stack.ParseArgs(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error")
		return
	}

	// Already sorted → print nothing
	if stack.IsSorted(a) {
		return
	}

	b := stack.Stack{}

	// Convert values to indexes (0..n-1)
	a = stack.Indexify(a)

	// Radix sort
	radixSort(&a, &b)
}

// ---------- RADIX SORT ----------

func radixSort(a, b *stack.Stack) {
	size := len(*a)
	maxBits := stack.MaxBits(*a)

	for bit := 0; bit < maxBits; bit++ {
		for i := 0; i < size; i++ {
			if ((*a)[0]>>bit)&1 == 0 {
				stack.Pb(a, b)
				fmt.Println("pb")
			} else {
				stack.Ra(a)
				fmt.Println("ra")
			}
		}

		for len(*b) > 0 {
			stack.Pa(a, b)
			fmt.Println("pa")
		}
	}
}
