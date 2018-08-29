package fizzbuzz

import (
	"errors"
	"fmt"
	"io"
)

/*
Results holds the starting number and an array of lines for each number
from the start up to the end number given to the Process function.
*/
type Results struct {
	Start, End int
	Lines      chan string
}

/*
Size returns the number of items in the result
*/
func (r *Results) Size() int {
	return r.End - r.Start
}

/*
Print writes the lines of the results to the writer
*/
func (r *Results) Print(w io.Writer) {
	i := 0
	for line := range r.Lines {
		fmt.Fprintf(w, "%v: %v\n", i+r.Start, line)
		i++
	}
}

/*
Process generates Results containing the following for each number from start (inclusive) to end (exclusive):
if the number is divisible by 3, a line with "Fizz"
if the number is divisible by 5, a line with "Buzz"
if the number is divisible by both, a line with "FizzBuzz"
blank lines for all other numbers
error is returned if start is greater than end
*/
func Process(start int, end int) (*Results, error) {
	if start > end {
		return nil, errors.New("start cannot be greater than end")
	}
	lines := make(chan string)
	rval := Results{start, end, lines}
	go generate(&rval)
	return &rval, nil
}

func generate(results *Results) {
	count := results.Size()
	num := results.Start
	for i := 0; i < count; i++ {
		results.Lines <- Eval(num)
		num++
	}
	close(results.Lines)
}

/*
Eval takes in an integer and returns
an empty string, Fizz if divisible by three, Buzz if divisible by five or FizzBuzz if both
*/
func Eval(i int) string {
	rval := ""
	if i%3 == 0 {
		rval += "Fizz"
	}
	if i%5 == 0 {
		rval += "Buzz"
	}
	return rval
}
