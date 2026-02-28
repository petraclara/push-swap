# Push-Swap & Checker

A Go implementation of the Push-Swap project: an efficient non-comparative sorting algorithm using two stacks and a limited set of operations.

## Project Logic

### The Stacks
The project uses two stacks, **A** and **B**:
- **Stack A**: Initially contains a random list of unique integers.
- **Stack B**: Initially empty, used as a workspace for sorting.

### Permitted Operations
- `pa`, `pb`: Push the top element from one stack to another.
- `sa`, `sb`, `ss`: Swap the first two elements of a stack.
- `ra`, `rb`, `rr`: Rotate all elements up by one (shift up).
- `rra`, `rrb`, `rrr`: Reverse rotate all elements down by one (shift down).

### Sorting Algorithms
The `push-swap` program chooses the most efficient strategy based on the input size:
1.  **Small Sets (n ≤ 3)**: Uses a hardcoded decision tree to sort in 2-3 moves.
2.  **Medium Sets (n ≤ 5)**: Finds the smallest elements, pushes them to B, sorts the remaining 3 in A, and pushes back.
3.  **Large Sets (n > 5)**: Implements **Radix Sort** after "Indexifying" the input (mapping values to their relative ranks 0 to n-1). This ensures efficiency even with negative or large numbers.

---

## How it Runs

### Compilation
You can compile both programs using the provided `Makefile`:
```bash
make
```
This will generate the `push-swap` and `checker` executables.

Other commands:
- `make clean`: Removes the executables.
- `make re`: Recompiles everything.
- `make test`: Runs the unit tests.

### Usage
Run `push-swap` to get the instructions, or pipe it directly into `checker` to verify the result.

```bash
# Get instructions only
./push-swap "4 67 3 87 23"

# Verify with checker (should print OK)
ARG="4 67 3 87 23"; ./push-swap "$ARG" | ./checker "$ARG"
```

### Testing
This project includes comprehensive unit tests for stack operations and utility functions.
```bash
go test -v ./stack/...
```

## Error Handling
The programs will display `Error` on standard error if:
- Arguments are not integers.
- There are duplicate integers.
- `checker` receives an invalid instruction.
