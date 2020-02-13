package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/bclement/codenight/pkg/postfix"
)

type source interface {
	Scan() bool
	Text() string
	Err() error
}

type readerSource struct {
	*bufio.Scanner
}

func newReaderSource(in io.Reader) readerSource {
	scanner := bufio.NewScanner(in)
	scanner.Split(bufio.ScanWords)
	return readerSource{scanner}
}

type sliceSource struct {
	s []string
	i int
}

func newSliceSource(s []string) *sliceSource {
	return &sliceSource{s: s, i: -1}
}

func (s *sliceSource) Scan() bool {
	s.i++
	return s.i < len(s.s)
}

func (s sliceSource) Text() string {
	return s.s[s.i]
}

func (s sliceSource) Err() error {
	return nil
}

func main() {
	var err error
	var src source
	if len(os.Args) == 1 {
		src = newReaderSource(os.Stdin)
	} else if len(os.Args) > 3 {
		src = newSliceSource(os.Args[1:])
	} else {
		fmt.Printf("Invalid arguments\n")
		return
	}
	calc := postfix.NewCalc()
	for src.Scan() {
		t := src.Text()
		f, err := strconv.ParseFloat(t, 64)
		if err != nil {
			err = calc.SubmitOperator(t)
			if err != nil {
				fmt.Printf("Invalid input: %v\n", err)
				return
			}
		} else {
			calc.SubmitNumber(f)
		}
	}
	err = src.Err()
	if err != nil {
		fmt.Printf("Problem reading input: %v\n", err)
		return
	}
	fmt.Printf("%v\n", calc.Result())
}
