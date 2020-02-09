package stack

import "testing"

func TestEmptyPop(t *testing.T) {
	s := NewStack()
	_, valuePresent := s.Pop()
	if valuePresent {
		t.Fatalf("Expected no value present when popping on empty")
	}
}

func TestLIFO(t *testing.T) {
	s := NewStack()
	assertLen(t, s, 0, "create")
	s.Push(42)
	assertLen(t, s, 1, "push1")
	s.Push(17)
	assertLen(t, s, 2, "push2")
	i, _ := s.Peek()
	assertLen(t, s, 2, "peek")
	assertVal(t, 17, i, "peek")
	i, _ = s.Pop()
	assertLen(t, s, 1, "pop1")
	assertVal(t, 17, i, "pop1")
	i, _ = s.Pop()
	assertLen(t, s, 0, "pop2")
	assertVal(t, 42, i, "pop2")
}

func assertLen(t *testing.T, s Stack, expected int, op string) {
	l := s.Size()
	if l != expected {
		t.Fatalf("Expected %v length after %v, got %v", expected, op, l)
	}
}

func assertVal(t *testing.T, expected, actual float64, op string) {
	if expected != actual {
		t.Fatalf("Expected %v from %v, got %v", expected, op, actual)
	}
}
