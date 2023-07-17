package main

import (
	"bytes"
	"compress/gzip"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
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
	zipNew, zipOld := dataCompress(newFile), dataCompress(oldFile)
	if zipNew == nil || zipOld == nil {
		fmt.Println("Data compressing error.")
		return
	}
	zipDataCompare(zipNew, zipOld)
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

func dataCompress(file *os.File) (data []byte) {
	var compressData bytes.Buffer
	gzipWritter := gzip.NewWriter(&compressData)
	defer gzipWritter.Close()
	_, err := io.Copy(gzipWritter, file)
	if err != nil {
		return nil
	}
	return compressData.Bytes()
}

func zipDataCompare(bytes1, bytes2 []byte) {
	lines1 := zipDataSplit(bytes1)
	lines2 := zipDataSplit(bytes2)
	map1 := make(map[string]bool, len(lines1))
	for _, line := range lines1 {
		map1[line] = false
	}
	for _, line := range lines2 {
		if _, ok := map1[line]; !ok {
			fmt.Printf("ADDED %s\n", line)
			continue
		}
		map1[line] = true
	}
	for _, line := range lines1 {
		if map1[line] == false {
			fmt.Printf("REMOVED %s\n", line)
		}
	}
}

func zipDataSplit(data []byte) []string {
	str := string(data)
	arr := strings.Split(str, "\n")
	return arr
}
