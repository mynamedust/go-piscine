package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	oldPath := flag.String("old", "", "Write valid .txt file name")
	newPath := flag.String("new", "", "Write valid .txt file name")
	flag.Parse()
	oldFile, newFile := checkFile(*oldPath), checkFile(*newPath)
	defer oldFile.Close()
	defer newFile.Close()
	if oldFile == nil || newFile == nil {
		return
	}
	dataCompare(oldFile, newFile)
}

func checkFile(path string) *os.File {
	ext := filepath.Ext(path)
	if ext != ".txt" {
		fmt.Println("Write valid .txt file name")
		return nil
	}
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("File open failed.")
		return nil
	}
	return file
}

func dataCompare(oldFile, newFile *os.File) {
	oldData := make(map[string]struct{})
	rOld := bufio.NewScanner(oldFile)
	rNew := bufio.NewScanner(newFile)
	for rOld.Scan() {
		oldData[rOld.Text()] = struct{}{}
	}
	for rNew.Scan() {
		if _, ok := oldData[rNew.Text()]; !ok {
			fmt.Printf("ADDED %s\n", rNew.Text())
			continue
		}
		delete(oldData, rNew.Text())
	}
	for key := range oldData {
		fmt.Printf("REMOVED %s\n", key)
	}
}
