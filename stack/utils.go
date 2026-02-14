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
			return nil, errors.New("Error")
		}
		if seen[n] {
			return nil, errors.New("Error")
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
	leng := int(copy(sorted, a))

	for i := 0; i < leng; i++ {
		for j := i + 1; j < leng; j++ {
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

func MaxBits(a Stack) int {
	max := 0
	for _, v := range a {
		if v > max {
			max = v
		}
	}

	bits := 0
	for (max >> bits) != 0 {
		bits++
	}
	return bits
}
