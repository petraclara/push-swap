package main

import (
	"fmt"
	"os"
	"push-swap/stack"
)

func main() {
	if len(os.Args) < 2 {
		return
	}

	a, err := stack.ParseArgs(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error")
		return
	}

	if stack.IsSorted(a) {
		return
	}

	b := stack.Stack{}
	a = stack.Indexify(a)

	size := len(a)

	switch {
	case size == 2:
		sortTwo(&a)
	case size == 3:
		sortThree(&a)
	case size <= 5:
		sortSmallTurk(&a, &b)
	case size <= 100:
		turkSort(&a, &b)
	default:
		radixSort(&a, &b)
	}
}

// ---------- SORT TWO ----------
func sortTwo(a *stack.Stack) {
	if (*a)[0] > (*a)[1] {
		stack.Sa(a)
		fmt.Println("sa")
	}
}

// ---------- SORT THREE ----------
func sortThree(a *stack.Stack) {
	first := (*a)[0]
	second := (*a)[1]
	third := (*a)[2]

	if first > second && second < third && first < third {
		stack.Sa(a)
		fmt.Println("sa")
	} else if first > second && second > third {
		stack.Sa(a)
		fmt.Println("sa")
		stack.Rra(a)
		fmt.Println("rra")
	} else if first > second && second < third && first > third {
		stack.Ra(a)
		fmt.Println("ra")
	} else if first < second && second > third && first < third {
		stack.Sa(a)
		fmt.Println("sa")
		stack.Ra(a)
		fmt.Println("ra")
	} else if first < second && second > third && first > third {
		stack.Rra(a)
		fmt.Println("rra")
	}
}

// ---------- TURK ALGORITHM ----------
func turkSort(a, b *stack.Stack) {
	size := len(*a)

	// Step 1: Push first 2 to B
	if size > 3 {
		stack.Pb(a, b)
		fmt.Println("pb")
	}
	if size > 4 {
		stack.Pb(a, b)
		fmt.Println("pb")
	}

	// Step 2: Push rest to B using cost calculation (leave 3 in A)
	for len(*a) > 3 {
		pushCheapestOptimized(a, b)
	}

	// Step 3: Sort remaining 3 in A
	if len(*a) == 3 {
		sortThree(a)
	} else if len(*a) == 2 {
		sortTwo(a)
	}

	// Step 4: Push everything back from B to A
	for len(*b) > 0 {
		pushBackToA(a, b)
	}
}

// Simpler algorithm for 4-5 elements
func sortSmallTurk(a, b *stack.Stack) {
	// Push all but 3 to B (push the smallest values)
	for len(*a) > 3 {
		minIdx := findMinIndex(a)
		aLen := len(*a)

		// Rotate min to top
		if minIdx <= aLen/2 {
			for i := 0; i < minIdx; i++ {
				stack.Ra(a)
				fmt.Println("ra")
			}
		} else {
			for i := 0; i < aLen-minIdx; i++ {
				stack.Rra(a)
				fmt.Println("rra")
			}
		}

		stack.Pb(a, b)
		fmt.Println("pb")
	}

	// Sort remaining 3
	if len(*a) == 3 {
		sortThree(a)
	} else if len(*a) == 2 {
		sortTwo(a)
	}

	// Push back in order (smallest first)
	for len(*b) > 0 {
		stack.Pa(a, b)
		fmt.Println("pa")
	}
}

// Helper: find index of minimum value
func findMinIndex(a *stack.Stack) int {
	minIdx := 0
	for i := 1; i < len(*a); i++ {
		if (*a)[i] < (*a)[minIdx] {
			minIdx = i
		}
	}
	return minIdx
}

// Optimized version that can use simultaneous rotations
func pushCheapestOptimized(a, b *stack.Stack) {
	minCost := int(^uint(0) >> 1)
	bestIdx := 0
	bestRotA := 0
	bestRotB := 0
	bestDirA := ""
	bestDirB := ""

	// Calculate cost for each element in A
	for i := 0; i < len(*a); i++ {
		value := (*a)[i]
		aLen := len(*a)

		// Calculate A rotations
		rotA := i
		dirA := "ra"
		if i > aLen/2 {
			rotA = aLen - i
			dirA = "rra"
		}

		// Calculate B rotations
		rotB := 0
		dirB := "rb"
		if len(*b) > 0 {
			targetIdx := findTargetInB(b, value)
			bLen := len(*b)
			rotB = targetIdx
			dirB = "rb"
			if targetIdx > bLen/2 {
				rotB = bLen - targetIdx
				dirB = "rrb"
			}
		}

		// Calculate total cost (can do simultaneous rotations)
		cost := rotA + rotB
		if (dirA == "ra" && dirB == "rb") || (dirA == "rra" && dirB == "rrb") {
			// Can rotate both simultaneously - cost is the max
			if rotA > rotB {
				cost = rotA
			} else {
				cost = rotB
			}
		}

		if cost < minCost {
			minCost = cost
			bestIdx = i
			bestRotA = rotA
			bestRotB = rotB
			bestDirA = dirA
			bestDirB = dirB
		}
	}

	// Execute the move with optimization
	executeMoveOptimized(a, b, bestIdx, bestRotA, bestRotB, bestDirA, bestDirB)
}

// Execute move with simultaneous rotation optimization
func executeMoveOptimized(a, b *stack.Stack, idx, rotA, rotB int, dirA, dirB string) {
	// Do simultaneous rotations if same direction
	if (dirA == "ra" && dirB == "rb") || (dirA == "rra" && dirB == "rrb") {
		minRot := rotA
		if rotB < rotA {
			minRot = rotB
		}

		for i := 0; i < minRot; i++ {
			if dirA == "ra" {
				stack.Rr(a, b)
				fmt.Println("rr")
			} else {
				stack.Rrr(a, b)
				fmt.Println("rrr")
			}
		}
		rotA -= minRot
		rotB -= minRot
	}

	// Finish A rotations
	for i := 0; i < rotA; i++ {
		if dirA == "ra" {
			stack.Ra(a)
			fmt.Println("ra")
		} else {
			stack.Rra(a)
			fmt.Println("rra")
		}
	}

	// Finish B rotations
	for i := 0; i < rotB; i++ {
		if dirB == "rb" {
			stack.Rb(b)
			fmt.Println("rb")
		} else {
			stack.Rrb(b)
			fmt.Println("rrb")
		}
	}

	// Push to B
	stack.Pb(a, b)
	fmt.Println("pb")
}

// Find where in B this value should be inserted
func findTargetInB(b *stack.Stack, value int) int {
	if len(*b) == 0 {
		return 0
	}

	// Find the position where value should go
	// B is sorted in descending order (largest on top)
	targetIdx := 0
	maxFound := false

	for i := 0; i < len(*b); i++ {
		if (*b)[i] < value {
			if !maxFound || (*b)[i] > (*b)[targetIdx] {
				targetIdx = i
				maxFound = true
			}
		}
	}

	// If value is smaller than all in B, put on top of largest
	if !maxFound {
		maxIdx := 0
		for i := 1; i < len(*b); i++ {
			if (*b)[i] > (*b)[maxIdx] {
				maxIdx = i
			}
		}
		targetIdx = maxIdx
	}

	return targetIdx
}

// Push back from B to A in correct order
func pushBackToA(a, b *stack.Stack) {
	// Find max in B
	maxIdx := 0
	for i := 1; i < len(*b); i++ {
		if (*b)[i] > (*b)[maxIdx] {
			maxIdx = i
		}
	}

	// Rotate to bring max to top
	bLen := len(*b)
	if maxIdx <= bLen/2 {
		for i := 0; i < maxIdx; i++ {
			stack.Rb(b)
			fmt.Println("rb")
		}
	} else {
		for i := 0; i < bLen-maxIdx; i++ {
			stack.Rrb(b)
			fmt.Println("rrb")
		}
	}

	// Push to A
	stack.Pa(a, b)
	fmt.Println("pa")
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
