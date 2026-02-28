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

	// Join all arguments into a single string
	argString := strings.Join(os.Args[1:], " ")

	// Parse initial stack A
	a, err := stack.ParseArgs(argString)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: failed to parse arguments on checker program")
		return
	}

	b := stack.Stack{}

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		instr := strings.TrimSpace(scanner.Text())
		if instr == "" {
			continue
		}

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
		default:
			fmt.Fprintln(os.Stderr, "Error: failed, encountered an invalid operation on checker")
			return
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Error: failed to scan")
		return
	}

	if stack.IsSorted(a) && len(b) == 0 {
		fmt.Println("OK")
	} else {
		fmt.Println("KO")
	}
}
