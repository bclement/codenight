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
	Start int
	Lines []string
}

/*
Print writes the lines of the results to the writer
*/
func (r *Results) Print(w io.Writer) {
	for i := 0; i < len(r.Lines); i++ {
		fmt.Fprintf(w, "%v: %v\n", i+r.Start, r.Lines[i])
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
	count := end - start
	num := start
	rval := Results{start, make([]string, count)}
	for i := 0; i < count; i++ {
		rval.Lines[i] = Eval(num)
		num++
	}
	return &rval, nil
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
