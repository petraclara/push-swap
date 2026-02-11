package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"push-swap/stack"
)

func main() {
	// No arguments → display nothing
	if len(os.Args) < 2 {
		return
	}

	// Parse initial stack A
	a, err := stack.ParseArgs(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error")
		return
	}
	b := stack.Stack{}

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		instr := strings.TrimSpace(scanner.Text())

		switch instr {
		case "sa":
			stack.Sa(&a)
		case "sb":
			stack.Sb(&b)
		case "ss":
			stack.Ss(&a, &b)
		case "pa":
			stack.Pa(&a, &b)
		case "pb":
			stack.Pb(&a, &b)
		case "ra":
			stack.Ra(&a)
		case "rb":
			stack.Rb(&b)
		case "rr":
			stack.Rr(&a, &b)
		case "rra":
			stack.Rra(&a)
		case "rrb":
			stack.Rrb(&b)
		case "rrr":
			stack.Rrr(&a, &b)
		case "":
			// ignore empty lines
		default:
			fmt.Fprintln(os.Stderr, "Error")
			return
		}
	}

	// Scanner error
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Error")
		return
	}

	// Final check
	if stack.IsSorted(a) && len(b) == 0 {
		fmt.Println("OK")
	} else {
		fmt.Println("KO")
	}
}
