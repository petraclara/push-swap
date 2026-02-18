package stack

import (
	"reflect"
	"testing"
)

func TestSa(t *testing.T) {
	s := Stack{1, 2, 3}
	Sa(&s)
	expected := Stack{2, 1, 3}
	if !reflect.DeepEqual(s, expected) {
		t.Errorf("Expected %v, got %v", expected, s)
	}

	s2 := Stack{1}
	Sa(&s2)
	if !reflect.DeepEqual(s2, Stack{1}) {
		t.Errorf("Expected {1}, got %v", s2)
	}
}

func TestPaPb(t *testing.T) {
	a := Stack{1, 2}
	b := Stack{3, 4}
	
	Pb(&a, &b)
	if !reflect.DeepEqual(a, Stack{2}) || !reflect.DeepEqual(b, Stack{1, 3, 4}) {
		t.Errorf("Pb failed: a=%v, b=%v", a, b)
	}

	Pa(&a, &b)
	if !reflect.DeepEqual(a, Stack{1, 2}) || !reflect.DeepEqual(b, Stack{3, 4}) {
		t.Errorf("Pa failed: a=%v, b=%v", a, b)
	}
}

func TestRa(t *testing.T) {
	s := Stack{1, 2, 3}
	Ra(&s)
	expected := Stack{2, 3, 1}
	if !reflect.DeepEqual(s, expected) {
		t.Errorf("Expected %v, got %v", expected, s)
	}
}

func TestRra(t *testing.T) {
	s := Stack{1, 2, 3}
	Rra(&s)
	expected := Stack{3, 1, 2}
	if !reflect.DeepEqual(s, expected) {
		t.Errorf("Expected %v, got %v", expected, s)
	}
}

func TestSsRrRrr(t *testing.T) {
	a := Stack{1, 2}
	b := Stack{3, 4}
	
	Ss(&a, &b)
	if !reflect.DeepEqual(a, Stack{2, 1}) || !reflect.DeepEqual(b, Stack{4, 3}) {
		t.Errorf("Ss failed: a=%v, b=%v", a, b)
	}

	Rr(&a, &b)
	if !reflect.DeepEqual(a, Stack{1, 2}) || !reflect.DeepEqual(b, Stack{3, 4}) {
		t.Errorf("Rr failed: a=%v, b=%v", a, b)
	}

	Rrr(&a, &b)
	if !reflect.DeepEqual(a, Stack{2, 1}) || !reflect.DeepEqual(b, Stack{4, 3}) {
		t.Errorf("Rrr failed: a=%v, b=%v", a, b)
	}
}
