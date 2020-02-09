package postfix

import (
	"fmt"
	"math"

	"github.com/bclement/codenight/pkg/stack"
)

/*
Calculator performs postfix calculations
*/
type Calculator struct {
	s stack.Stack
}

/*
NewCalc creates a new calculator
*/
func NewCalc() *Calculator {
	return &Calculator{stack.NewStack()}
}

/*
SubmitNumber takes the next digit in the calculation
*/
func (c *Calculator) SubmitNumber(i float64) {
	c.s.Push(i)
}

/*
Result shows the current total
*/
func (c *Calculator) Result() float64 {
	i, ok := c.s.Peek()
	if !ok {
		return 0
	}
	return i
}

/*
Reset clears the calculator
*/
func (c *Calculator) Reset() {
	c.s.Reset()
}

/*
SubmitOperator submits the next operator for the calculation
returns an error if the operator is not recognized or
calculator cannot handle operation in current state
*/
func (c *Calculator) SubmitOperator(o string) error {
	if c.s.Size() < 2 {
		return fmt.Errorf("Not enough values on stack for operator: %v", o)
	}
	var f func(float64, float64) float64
	switch o {
	case "+":
		f = add
	case "-":
		f = sub
	case "*":
		f = mult
	case "/":
		f = div
	case "^":
		f = pow
	default:
		return fmt.Errorf("Unknown operator: %v", o)
	}
	i, _ := c.s.Pop()
	j, _ := c.s.Pop()
	res := f(i, j)
	c.s.Push(res)
	return nil
}

func add(i, j float64) float64 {
	return i + j
}

func sub(i, j float64) float64 {
	return i - j
}

func mult(i, j float64) float64 {
	return i * j
}

func div(i, j float64) float64 {
	return i / j
}

func pow(i, j float64) float64 {
	return math.Pow(i, j)
}
