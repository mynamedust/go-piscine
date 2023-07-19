package main

import (
	"flag"
	"fmt"
	"path/filepath"
	"sync"
)

func main() {
	wFlag := flag.Bool("w", false, "Use -w flag for count words.")
	lFlag := flag.Bool("l", false, "Use -w flag for count lines.")
	mFlag := flag.Bool("m", false, "Use -w flag for count characters.")
	flag.Parse()
	args := flag.Args()
	if (*wFlag && *lFlag) || (*wFlag && *mFlag) || (*lFlag && *mFlag) {
		fmt.Println("Error. There must be only one active flag.")
		return
	}
	if len(args) <= 0 {
		fmt.Println("Error. There must be at least one argument.")
		return
	}
	counter := selectCounter(*wFlag, *lFlag, *mFlag)
	countItems(counter, args)
}

func selectCounter(w, l, m bool) func(args string, wg *sync.WaitGroup) int {
	switch {
	case w:
		return wordCounter
	case l:
		return lineCounter
	case m:
		return charCounter
	default:
		return wordCounter
	}
}

func countItems(counter func(string, *sync.WaitGroup) int, args []string) {
	wg := new(sync.WaitGroup)
	for _, name := range args {
		if !validName(name) {
			continue
		}
		wg.Add(1)
		go counter(name, wg)
	}
	wg.Wait()
}

func validName(name string) bool {
	ext := filepath.Ext(name)
	if ext != ".txt" {
		fmt.Println("Error.Wrong file type.")
		return false
	}
	return true
}

func wordCounter(file string, wg *sync.WaitGroup) (count int) {
	defer wg.Done()
	fmt.Println(file)
}

func lineCounter(file string, wg *sync.WaitGroup) (count int) {
	defer wg.Done()
	fmt.Println(file)
}

func charCounter(file string, wg *sync.WaitGroup) (count int) {
	defer wg.Done()
	fmt.Println(file)
}
