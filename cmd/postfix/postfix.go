package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/bclement/codenight/pkg/postfix"
)

func read(in io.Reader) (tokens []string, err error) {
	s := bufio.NewScanner(in)
	s.Split(bufio.ScanWords)
	for s.Scan() {
		tokens = append(tokens, s.Text())
	}
	err = s.Err()
	return
}

func main() {
	var err error
	var tokens []string
	if len(os.Args) == 1 {
		tokens, err = read(os.Stdin)
		if err != nil {
			fmt.Printf("Invalid input: %v\n", err)
		}
	} else if len(os.Args) > 3 {
		tokens = os.Args[1:]
	} else {
		fmt.Printf("Invalid arguments\n")
		return
	}
	calc := postfix.NewCalc()
	for _, t := range tokens {
		f, err := strconv.ParseFloat(t, 64)
		if err != nil {
			err = calc.SubmitOperator(t)
			if err != nil {
				fmt.Printf("Invalid input: %v\n", err)
			}
		} else {
			calc.SubmitNumber(f)
		}
	}
	fmt.Printf("%v\n", calc.Result())
}
