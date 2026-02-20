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
		fmt.Println(err)
		return
	}

	// Already sorted → print nothing
	if stack.IsSorted(a) {
		return
	}

	b := stack.Stack{}

	// Convert values to indexes (0..n-1)
	a = stack.Indexify(a)

	if len(a) <= 3 {
		sortSmall3(&a)
	} else if len(a) <= 12 {
		sortSmall12(&a, &b)
	} else {
		chunkSort(&a, &b)
	}
}

// ---------- SMALL SET SORT ----------

func sortSmall3(a *stack.Stack) {
	if stack.IsSorted(*a) {
		return
	}
	if len(*a) == 2 {
		if (*a)[0] > (*a)[1] {
			stack.Sa(a)
			fmt.Println("sa")
		}
		return
	}

	v0, v1, v2 := (*a)[0], (*a)[1], (*a)[2]

	if v0 < v1 && v1 > v2 && v0 < v2 {
		// 0 2 1
		stack.Sa(a)
		fmt.Println("sa")
		stack.Ra(a)
		fmt.Println("ra")
	} else if v0 > v1 && v1 < v2 && v0 < v2 {
		// 1 0 2
		stack.Sa(a)
		fmt.Println("sa")
	} else if v0 < v1 && v1 > v2 && v0 > v2 {
		// 1 2 0
		stack.Rra(a)
		fmt.Println("rra")
	} else if v0 > v1 && v1 < v2 && v0 > v2 {
		// 2 0 1
		stack.Ra(a)
		fmt.Println("ra")
	} else if v0 > v1 && v1 > v2 {
		// 2 1 0
		stack.Sa(a)
		fmt.Println("sa")
		stack.Rra(a)
		fmt.Println("rra")
	}
}

func sortSmall12(a, b *stack.Stack) {
	if stack.IsSorted(*a) {
		return
	}

	// Simple sa optimization
	if len(*a) > 1 && (*a)[0] > (*a)[1] {
		stack.Sa(a)
		fmt.Println("sa")
	}

	for len(*a) > 3 {
		if stack.IsSorted(*a) {
			break
		}
		minIdx := 0
		for i, v := range *a {
			if v < (*a)[minIdx] {
				minIdx = i
			}
		}

		if minIdx <= len(*a)/2 {
			for i := 0; i < minIdx; i++ {
				stack.Ra(a)
				fmt.Println("ra")
			}
		} else {
			for i := 0; i < len(*a)-minIdx; i++ {
				stack.Rra(a)
				fmt.Println("rra")
			}
		}
		stack.Pb(a, b)
		fmt.Println("pb")
	}

	sortSmall3(a)

	for len(*b) > 0 {
		stack.Pa(a, b)
		fmt.Println("pa")
	}
}

// ---------- CHUNK SORT ----------

func chunkSort(a, b *stack.Stack) {
	size := len(*a)
	chunkSize := 15
	if size > 100 {
		chunkSize = 35
	}

	for i := 0; len(*a) > 0; {
		if (*a)[0] <= i {
			stack.Pb(a, b)
			fmt.Println("pb")
			stack.Rb(b)
			fmt.Println("rb")
			i++
		} else if (*a)[0] <= i+chunkSize {
			stack.Pb(a, b)
			fmt.Println("pb")
			i++
		} else {
			stack.Ra(a)
			fmt.Println("ra")
		}
	}

	for len(*b) > 0 {
		maxIdx := 0
		for i, v := range *b {
			if v > (*b)[maxIdx] {
				maxIdx = i
			}
		}

		if maxIdx <= len(*b)/2 {
			for i := 0; i < maxIdx; i++ {
				stack.Rb(b)
				fmt.Println("rb")
			}
		} else {
			for i := 0; i < len(*b)-maxIdx; i++ {
				stack.Rrb(b)
				fmt.Println("rrb")
			}
		}
		stack.Pa(a, b)
		fmt.Println("pa")
	}
}
