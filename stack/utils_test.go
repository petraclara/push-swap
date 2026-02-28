package stack

import (
	"reflect"
	"testing"
)

func TestParseArgs(t *testing.T) {
	tests := []struct {
		input    string
		expected Stack
		wantErr  bool
	}{
		{"1 2 3", Stack{1, 2, 3}, false},
		{" 10   20  30 ", Stack{10, 20, 30}, false},
		{"", Stack{}, false},
		{"1 2 2", nil, true},
		{"1 two 3", nil, true},
	}

	for _, tt := range tests {
		got, err := ParseArgs(tt.input)
		if (err != nil) != tt.wantErr {
			t.Errorf("ParseArgs(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			continue
		}
		if !tt.wantErr && !reflect.DeepEqual(got, tt.expected) {
			t.Errorf("ParseArgs(%q) = %v, want %v", tt.input, got, tt.expected)
		}
	}
}

func TestIsSorted(t *testing.T) {
	if !IsSorted(Stack{1, 2, 3}) {
		t.Error("Expected {1, 2, 3} to be sorted")
	}
	if IsSorted(Stack{3, 2, 1}) {
		t.Error("Did not expect {3, 2, 1} to be sorted")
	}
}

func TestIndexify(t *testing.T) {
	s := Stack{100, 50, 150}
	got := Indexify(s)
	expected := Stack{1, 0, 2}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %v, got %v", expected, got)
	}
}

func TestMaxBits(t *testing.T) {
	if MaxBits(Stack{0, 1, 2, 3}) != 2 {
		t.Errorf("Expected 2 bits for 3, got %d", MaxBits(Stack{0, 1, 2, 3}))
	}
}
