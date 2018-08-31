package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/bclement/codenight/pkg/fizzbuzz"
)

func handler(w http.ResponseWriter, r *http.Request) {
	nums, err := getNums(r)
	if err != nil {
		msg := fmt.Sprintf("Problem parsing path as integer: %v", err)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}
	count := len(nums)
	if count == 0 {
		fmt.Fprintf(w, "Welcome to FizzBuzz as a Service!")
	} else if count == 1 {
		fmt.Fprintf(w, fizzbuzz.Eval(nums[0]))
	} else if count == 2 {
		results, err := fizzbuzz.Process(nums[0], nums[1])
		if err != nil {
			msg := fmt.Sprintf("Invalid input: %v", err)
			http.Error(w, msg, http.StatusBadRequest)
			return
		}
		results.Print(w)
	} else {
		http.Error(w, "No idea what you are looking for", http.StatusNotFound)
		return
	}
}

func getNums(r *http.Request) ([]int, error) {
	path := r.URL.Path[1:]
	if path == "" {
		return nil, nil
	}
	parts := strings.Split(r.URL.Path[1:], "/")
	count := len(parts)
	rval := make([]int, count)
	for i := 0; i < count; i++ {
		num, err := strconv.Atoi(parts[i])
		if err != nil {
			return rval, err
		}
		rval[i] = num
	}
	return rval, nil
}

func main() {
	port := 8080
	var err error
	if len(os.Args) == 2 {
		port, err = strconv.Atoi(os.Args[1])
		if err != nil {
			log.Fatalf("Argument must be a valid port number: %v", err)
		}
	}
	http.HandleFunc("/", handler)
	addr := fmt.Sprintf(":%v", port)
	log.Fatal(http.ListenAndServe(addr, nil))
}
