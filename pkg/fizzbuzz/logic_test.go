package fizzbuzz

import (
	"testing"
)

func TestTwenty(t *testing.T) {
	expected := [20]string{}
	for _, index := range []int{0, 15} {
		expected[index] = "FizzBuzz"
	}
	for _, index := range []int{3, 6, 9, 12, 18} {
		expected[index] = "Fizz"
	}
	for _, index := range []int{5, 10} {
		expected[index] = "Buzz"
	}
	start := 0
	end := 20
	results, err := Process(start, end)
	if err != nil {
		t.Errorf("Error during processing: %v", err)
		return
	}
	if results.Start != 0 {
		t.Errorf("Results start at %v instead of %v", results.Start, start)
		return
	}
	actual := results.Lines
	if len(actual) != end {
		t.Errorf("Results length %v instead of %v", len(actual), end)
		return
	}
	for i := start; i < end; i++ {
		if expected[i] != actual[i] {
			t.Errorf("Expected %v at index %v, got %v", expected[i], i, actual[i])
		}
	}
}
