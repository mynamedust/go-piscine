package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
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

func selectCounter(w, l, m bool) func(args string, wg *sync.WaitGroup) {
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

func countItems(counter func(string, *sync.WaitGroup), args []string) {
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
		fmt.Printf("Error reading file \"%s\".Wrong file type.\n", name)
		return false
	}
	return true
}

func wordCounter(fileName string, wg *sync.WaitGroup) {
	var count int
	defer wg.Done()
	data := openRead(fileName)
	if data == nil {
		return
	}
	words := strings.Split(string(data), " ")
	fmt.Println(len(words))
	for _, elem := range words {
		if strings.Compare(elem, " ") != 0 {
			fmt.Println(elem, "!=", elem != " ")
			count++
			count += strings.Count(elem, "\n")
		}
	}
	fmt.Println(count)
}

func lineCounter(fileName string, wg *sync.WaitGroup) {
	var count int
	defer wg.Done()
	data := openRead(fileName)
	if data == nil {
		return
	}
	if len(data) > 0 {
		count = 1
	}
	for _, elem := range data {
		if elem == '\n' {
			count++
		}
	}
	fmt.Println(count)
}

func charCounter(fileName string, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println(fileName)
	return
}

func openRead(fileName string) []byte {
	file, err := os.Open(fileName)
	defer file.Close()
	if err != nil {
		fmt.Printf("File \"%s\" opening error: %s", fileName, err)
		return nil
	}
	data, err := io.ReadAll(file)
	if err != nil {
		fmt.Printf("File \"%s\" reading error: %s", fileName, err)
		return nil
	}
	return data
}
