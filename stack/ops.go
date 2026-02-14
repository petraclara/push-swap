package stack

type Stack []int

func Pa(a, b *Stack) {
	if len(*b) == 0 {
		return
	}
	val := (*b)[0]
	*b = (*b)[1:]
	*a = append(Stack{val}, (*a)...)
}

// func Pb(a, b *Stack){
// 	if len(*a) == 0{
// 		return
// 	}
// 	val := (*a)[0]
// 	*b = (*a)[1:]
// 	*a = append(Stack{val}, (*b)...)
// }

func Pb(a, b *Stack) {
	if len(*a) == 0 {
		return
	}
	val := (*a)[0]
	*a = (*a)[1:]
	*b = append(Stack{val}, (*b)...)
}

func Sa(a *Stack) {
	if len(*a) < 2 {
		return
	}
	(*a)[0], (*a)[1] = (*a)[1], (*a)[0]
}

func Sb(b *Stack) {
	if len(*b) < 2 {
		return
	}
	(*b)[0], (*b)[1] = (*b)[1], (*b)[0]
}

func Ss(a, b *Stack) {
	Sa(a)
	Sb(b)
}

func Ra(a *Stack) {
	if len(*a) < 2 {
		return
	}
	first := (*a)[0]
	*a = append((*a)[1:], first)
}

func Rb(b *Stack) {
	if len(*b) < 2 {
		return
	}
	first := (*b)[0]
	*b = append((*b)[1:], first)
}

func Rr(a, b *Stack) {
	Ra(a)
	Rb(b)
}

func Rra(a *Stack) {
	if len(*a) < 2 {
		return
	}
	last := (*a)[len(*a)-1]
	*a = append(Stack{last}, (*a)[:len(*a)-1]...)
}

func Rrb(b *Stack) {
	if len(*b) < 2 {
		return
	}
	last := (*b)[len(*b)-1]
	*b = append(Stack{last}, (*b)[:len(*b)-1]...)
}

func Rrr(a, b *Stack) {
	Rra(a)
	Rrb(b)
}
