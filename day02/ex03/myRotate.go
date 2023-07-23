package main

import (
	"flag"
	"fmt"
)

func main() {
	dir := flag.String("a", ".", "Write valid path for archives.")
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		fmt.Println("Error. There must be at least one argument.")
		return
	}
	createArchives(*dir, args)
}