package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/bclement/codenight/pkg/fizzbuzz"
)

func main() {
	// os.Args includes the name of the executable
	if len(os.Args) == 2 {
		os.Stdout.WriteString(fizzbuzz.Eval(getArgs(1)[0]))
		os.Stdout.WriteString("\n")
	} else if len(os.Args) == 3 {
		values := getArgs(2)
		results, err := fizzbuzz.Process(values[0], values[1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error during processing: %v\n", err)
		} else {
			results.Print(os.Stdout)
		}
	} else {
		log.Fatalln("Must provide one or two integers")
	}
}

func getArgs(count int) []int {
	rval := make([]int, count)
	for i := 0; i < count; i++ {
		num, err := strconv.Atoi(os.Args[i+1])
		if err != nil {
			log.Fatalf("Invalid integer %v: %v", os.Args[i], err)
		}
		rval[i] = num
	}
	return rval
}
