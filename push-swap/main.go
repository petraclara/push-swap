package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"push-swap/stack"
)

var operations = []string{}

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

	// Save original numbers before indexifying (for visualizer)
	original := make([]string, len(a))
	for i, v := range a {
		original[i] = strconv.Itoa(v)
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

	// Write ops.txt for visualizer
	opsFile, err := os.OpenFile("visualizer/ops.txt", os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer opsFile.Close()
	for _, v := range operations {
		_, err = opsFile.WriteString(v + "\n")
		if err != nil {
			log.Fatal(err.Error())
		}
	}

	// Write initial.txt for visualizer (original unsorted numbers)
	initFile, err := os.OpenFile("visualizer/initial.txt", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer initFile.Close()
	_, err = initFile.WriteString(strings.Join(original, " ") + "\n")
	if err != nil {
		log.Fatal(err.Error())
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
			operations = append(operations, "sa")
		}
		return
	}

	v0, v1, v2 := (*a)[0], (*a)[1], (*a)[2]

	if v0 < v1 && v1 > v2 && v0 < v2 {
		// 0 2 1
		stack.Sa(a)
		fmt.Println("sa")
		operations = append(operations, "sa")

		stack.Ra(a)
		fmt.Println("ra")
		operations = append(operations, "ra")

	} else if v0 > v1 && v1 < v2 && v0 < v2 {
		// 1 0 2
		stack.Sa(a)
		fmt.Println("sa")
		operations = append(operations, "sa")
	} else if v0 < v1 && v1 > v2 && v0 > v2 {
		// 1 2 0
		stack.Rra(a)
		fmt.Println("rra")
		operations = append(operations, "rra")
	} else if v0 > v1 && v1 < v2 && v0 > v2 {
		// 2 0 1
		stack.Ra(a)
		fmt.Println("ra")
		operations = append(operations, "ra")
	} else if v0 > v1 && v1 > v2 {
		// 2 1 0
		stack.Sa(a)
		fmt.Println("sa")
		operations = append(operations, "sa")
		stack.Rra(a)
		fmt.Println("rra")
		operations = append(operations, "rra")
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
		operations = append(operations, "sa")
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
				operations = append(operations, "ra")
			}
		} else {
			for i := 0; i < len(*a)-minIdx; i++ {
				stack.Rra(a)
				fmt.Println("rra")
				operations = append(operations, "rra")
			}
		}
		stack.Pb(a, b)
		fmt.Println("pb")
		operations = append(operations, "pb")
	}

	sortSmall3(a)

	for len(*b) > 0 {
		stack.Pa(a, b)
		fmt.Println("pa")
		operations = append(operations, "pa")
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
			operations = append(operations, "pb")
			stack.Rb(b)
			fmt.Println("rb")
			operations = append(operations, "rb")
			i++
		} else if (*a)[0] <= i+chunkSize {
			stack.Pb(a, b)
			fmt.Println("pb")
			operations = append(operations, "pb")
			i++
		} else {
			stack.Ra(a)
			fmt.Println("ra")
			operations = append(operations, "ra")
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
				operations = append(operations, "rb")
			}
		} else {
			for i := 0; i < len(*b)-maxIdx; i++ {
				stack.Rrb(b)
				fmt.Println("rrb")
				operations = append(operations, "rrb")
			}
		}
		stack.Pa(a, b)
		fmt.Println("pa")
		operations = append(operations, "pa")
	}
}
