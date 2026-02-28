package stack

import (
	"errors"
	"strconv"
	"strings"
)

func ParseArgs(arg string) (Stack, error) {
	if strings.TrimSpace(arg) == "" {
		return Stack{}, nil
	}

	parts := strings.Fields(arg)
	stack := Stack{}
	seen := make(map[int]bool)

	for _, p := range parts {
		n, err := strconv.Atoi(p)
		if err != nil {
			return nil, errors.New("Error: failed to run strconv.Atoi")
		}
		if seen[n] {
			return nil, errors.New("Error: failed, seen n")
		}
		seen[n] = true
		stack = append(stack, n)
	}

	return stack, nil

}

func IsSorted(a Stack) bool {
	for i := 0; i < len(a)-1; i++ {
		if a[i] > a[i+1] {
			return false
		}
	}
	return true
}

func Indexify(a Stack) Stack {
	sorted := make(Stack, len(a))
	copy(sorted, a)

	for i := 0; i < len(sorted)-1; i++ {
		for j := i + 1; j < len(sorted); j++ {
			if sorted[i] > sorted[j] {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}

	index := make(map[int]int)
	for i, v := range sorted {
		index[v] = i
	}

	result := make(Stack, len(a))
	for i, v := range a {
		result[i] = index[v]
	}

	return result
}
