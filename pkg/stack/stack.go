package stack

/*
Stack is a basic number stack interface
*/
type Stack interface {
	Push(float64)
	Peek() (float64, bool)
	Pop() (float64, bool)
	Size() int
	Reset()
}

type numstack struct {
	storage []float64
}

/*
NewStack returns a newly allocated stack implementation
*/
func NewStack() Stack {
	return &numstack{nil}
}

func (ns *numstack) Push(i float64) {
	ns.storage = append(ns.storage, i)
}

func (ns *numstack) Peek() (i float64, ok bool) {
	l := len(ns.storage)
	if l > 0 {
		i, ok = ns.storage[l-1], true
	}
	return
}

func (ns *numstack) Pop() (i float64, ok bool) {
	i, ok = ns.Peek()
	if ok {
		ns.storage = ns.storage[:len(ns.storage)-1]
	}
	return
}

func (ns *numstack) Size() int {
	return len(ns.storage)
}

func (ns *numstack) Reset() {
	ns.storage = nil
}
