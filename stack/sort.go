package stack

import "fmt"

func SortStack(a, b *Stack) {
	switch len(*a) {
	case 0, 1:
	case 2:
		if (*a)[0] > (*a)[1] {
			Sa(a)
			fmt.Println("sa")
		}
	default:
		if len(*a) <= 6 {
			solveSmall(a, b)
		} else {
			sortTurk(a, b)
		}
	}
}

// ── SMALL SORT: BFS finds optimal solution for n≤6 in <5000 states ──

type sstate struct{ a, b []int }

var allOps = []string{"sa", "sb", "ss", "pa", "pb", "ra", "rb", "rr", "rra", "rrb", "rrr"}

func stateKey(s sstate) string {
	k := make([]byte, 0, 16)
	for _, v := range s.a {
		k = append(k, byte(v+'0'))
	}
	k = append(k, '|')
	for _, v := range s.b {
		k = append(k, byte(v+'0'))
	}
	return string(k)
}

func solveSmall(a, b *Stack) {
	start := sstate{append([]int{}, *a...), nil}
	goal := sortedCopy(start.a)

	type node struct {
		s    sstate
		path []string
	}
	visited := map[string]bool{stateKey(start): true}
	queue := []node{{start, nil}}

	for qi := 0; qi < len(queue); qi++ {
		cur := queue[qi]
		for _, op := range allOps {
			ns := applyOp(cur.s, op)
			if intsEq(ns.a, goal) && len(ns.b) == 0 {
				for _, o := range cur.path {
					execOp(a, b, o)
				}
				execOp(a, b, op)
				return
			}
			sk := stateKey(ns)
			if !visited[sk] {
				visited[sk] = true
				newPath := make([]string, len(cur.path)+1)
				copy(newPath, cur.path)
				newPath[len(cur.path)] = op
				queue = append(queue, node{ns, newPath})
			}
		}
	}
}

func applyOp(s sstate, op string) sstate {
	a := append([]int{}, s.a...)
	b := append([]int{}, s.b...)
	switch op {
	case "sa":
		if len(a) >= 2 {
			a[0], a[1] = a[1], a[0]
		}
	case "sb":
		if len(b) >= 2 {
			b[0], b[1] = b[1], b[0]
		}
	case "ss":
		if len(a) >= 2 {
			a[0], a[1] = a[1], a[0]
		}
		if len(b) >= 2 {
			b[0], b[1] = b[1], b[0]
		}
	case "pa":
		if len(b) > 0 {
			a = append([]int{b[0]}, a...)
			b = b[1:]
		}
	case "pb":
		if len(a) > 0 {
			b = append([]int{a[0]}, b...)
			a = a[1:]
		}
	case "ra":
		if len(a) >= 2 {
			a = append(a[1:], a[0])
		}
	case "rb":
		if len(b) >= 2 {
			b = append(b[1:], b[0])
		}
	case "rr":
		if len(a) >= 2 {
			a = append(a[1:], a[0])
		}
		if len(b) >= 2 {
			b = append(b[1:], b[0])
		}
	case "rra":
		if len(a) >= 2 {
			a = append([]int{a[len(a)-1]}, a[:len(a)-1]...)
		}
	case "rrb":
		if len(b) >= 2 {
			b = append([]int{b[len(b)-1]}, b[:len(b)-1]...)
		}
	case "rrr":
		if len(a) >= 2 {
			a = append([]int{a[len(a)-1]}, a[:len(a)-1]...)
		}
		if len(b) >= 2 {
			b = append([]int{b[len(b)-1]}, b[:len(b)-1]...)
		}
	}
	return sstate{a, b}
}

func execOp(a, b *Stack, op string) {
	switch op {
	case "sa":
		Sa(a)
		fmt.Println("sa")
	case "sb":
		Sb(b)
		fmt.Println("sb")
	case "ss":
		Ss(a, b)
		fmt.Println("ss")
	case "pa":
		Pa(a, b)
		fmt.Println("pa")
	case "pb":
		Pb(a, b)
		fmt.Println("pb")
	case "ra":
		Ra(a)
		fmt.Println("ra")
	case "rb":
		Rb(b)
		fmt.Println("rb")
	case "rr":
		Rr(a, b)
		fmt.Println("rr")
	case "rra":
		Rra(a)
		fmt.Println("rra")
	case "rrb":
		Rrb(b)
		fmt.Println("rrb")
	case "rrr":
		Rrr(a, b)
		fmt.Println("rrr")
	}
}

// ── TURK: greedy cost-minimised sort for n≥7 ─────────────────────

func sortTurk(a, b *Stack) {
	for len(*a) > 3 {
		cheapPushToB(a, b)
	}
	sortThree(a)
	for len(*b) > 0 {
		pushBtoA(a, b)
	}
	rotTop(a, minIdx(*a))
}

func sortThree(a *Stack) {
	t, m, b := (*a)[0], (*a)[1], (*a)[2]
	switch {
	case t > m && m < b && t < b:
		Sa(a)
		fmt.Println("sa")
	case t < m && m > b && t < b:
		Sa(a)
		fmt.Println("sa")
		Ra(a)
		fmt.Println("ra")
	case t > m && m < b && t > b:
		Ra(a)
		fmt.Println("ra")
	case t < m && m > b && t > b:
		Rra(a)
		fmt.Println("rra")
	case t > m && m > b:
		Sa(a)
		fmt.Println("sa")
		Rra(a)
		fmt.Println("rra")
	}
}

func cheapPushToB(a, b *Stack) {
	type cand struct {
		ia, ib   int
		raA, raB bool
		cost     int
	}
	res := cand{cost: 1<<31 - 1}
	for i, val := range *a {
		ib := targetB(*b, val)
		cA, cB := rotCost(i, len(*a)), rotCost(ib, len(*b))
		raA, raB := i <= len(*a)/2, ib <= len(*b)/2
		c := cA + cB
		if raA == raB {
			c = max2(cA, cB)
		}
		c++
		if c < res.cost {
			res = cand{i, ib, raA, raB, c}
		}
	}
	doRotPush(a, b, res.ia, res.ib, res.raA, res.raB)
}

func doRotPush(a, b *Stack, ia, ib int, raA, raB bool) {
	remA, remB := ia, ib
	if !raA {
		remA = len(*a) - ia
	}
	if !raB {
		remB = len(*b) - ib
	}
	if raA && raB {
		for ; remA > 0 && remB > 0; remA, remB = remA-1, remB-1 {
			Rr(a, b)
			fmt.Println("rr")
		}
	}
	if !raA && !raB {
		for ; remA > 0 && remB > 0; remA, remB = remA-1, remB-1 {
			Rrr(a, b)
			fmt.Println("rrr")
		}
	}
	for ; remA > 0; remA-- {
		if raA {
			Ra(a)
			fmt.Println("ra")
		} else {
			Rra(a)
			fmt.Println("rra")
		}
	}
	for ; remB > 0; remB-- {
		if raB {
			Rb(b)
			fmt.Println("rb")
		} else {
			Rrb(b)
			fmt.Println("rrb")
		}
	}
	Pb(a, b)
	fmt.Println("pb")
}

func pushBtoA(a, b *Stack) { rotTop(a, targetA(*a, (*b)[0])); Pa(a, b); fmt.Println("pa") }

func targetB(b Stack, val int) int {
	if len(b) == 0 {
		return 0
	}
	best, bestV := -1, -(1 << 30)
	for i, v := range b {
		if v < val && v > bestV {
			bestV = v
			best = i
		}
	}
	if best == -1 {
		return maxIdx(b)
	}
	return best
}

func targetA(a Stack, val int) int {
	best, bestV := -1, 1<<31-1
	for i, v := range a {
		if v > val && v < bestV {
			bestV = v
			best = i
		}
	}
	if best == -1 {
		return minIdx(a)
	}
	return best
}

// ── HELPERS ──────────────────────────────────────────────────────

func rotTop(s *Stack, idx int) {
	n := len(*s)
	if idx <= n/2 {
		for i := 0; i < idx; i++ {
			Ra(s)
			fmt.Println("ra")
		}
	} else {
		for i := 0; i < n-idx; i++ {
			Rra(s)
			fmt.Println("rra")
		}
	}
}

func rotCost(idx, n int) int {
	if idx <= n/2 {
		return idx
	}
	return n - idx
}
func max2(a, b int) int {
	if a > b {
		return a
	}
	return b
}
func minIdx(s Stack) int {
	m, i := s[0], 0
	for j, v := range s {
		if v < m {
			m = v
			i = j
		}
	}
	return i
}
func maxIdx(s Stack) int {
	m, i := s[0], 0
	for j, v := range s {
		if v > m {
			m = v
			i = j
		}
	}
	return i
}

func sortedCopy(s []int) []int {
	c := append([]int{}, s...)
	for i := range c {
		for j := i + 1; j < len(c); j++ {
			if c[i] > c[j] {
				c[i], c[j] = c[j], c[i]
			}
		}
	}
	return c
}
func intsEq(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
