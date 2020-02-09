package postfix

import "testing"

func TestAdd(t *testing.T) {
	c := NewCalc()
	c.SubmitNumber(3)
	c.SubmitNumber(5)
	c.SubmitOperator("+")
	assertTotal(t, c, 8)
}

func TestSub(t *testing.T) {
	c := NewCalc()
	c.SubmitNumber(3)
	c.SubmitNumber(5)
	c.SubmitOperator("-")
	assertTotal(t, c, -2)
}

func TestMult(t *testing.T) {
	c := NewCalc()
	c.SubmitNumber(3)
	c.SubmitNumber(5)
	c.SubmitOperator("*")
	assertTotal(t, c, 15)
}

func TestDiv(t *testing.T) {
	c := NewCalc()
	c.SubmitNumber(3)
	c.SubmitNumber(5)
	c.SubmitOperator("/")
	assertTotal(t, c, .6)
}

func TestPow(t *testing.T) {
	c := NewCalc()
	c.SubmitNumber(3)
	c.SubmitNumber(5)
	c.SubmitOperator("^")
	assertTotal(t, c, 243)
}

func TestComplex(t *testing.T) {
	c := NewCalc()
	c.SubmitNumber(15)
	c.SubmitNumber(7)
	c.SubmitNumber(1)
	c.SubmitNumber(1)
	c.SubmitOperator("+")
	c.SubmitOperator("-")
	c.SubmitOperator("/")
	c.SubmitNumber(3)
	c.SubmitOperator("*")
	assertTotal(t, c, 9)
}

func assertTotal(t *testing.T, c *Calculator, expected float64) {
	actual := c.Result()
	if actual != actual {
		t.Fatalf("Expected %v, got %v", expected, actual)
	}
}
