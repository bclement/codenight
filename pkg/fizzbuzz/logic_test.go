package fizzbuzz

import (
	"fmt"
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
	if results.Size() != end {
		t.Errorf("Results length %v instead of %v", results.Size(), end)
		return
	}
	i := start
	for actual := range results.Lines {
		value := expected[i]
		if value == "" {
			value = fmt.Sprint(i)
		}
		if value != actual {
			t.Errorf("Expected %v at index %v, got %v", value, i, actual)
		}
		i++
	}
}
