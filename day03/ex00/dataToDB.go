package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
)

func main() {
	file, err := os.Open("../materials/data.csv")
	if err != nil {
		fmt.Println("error: file opening failed: ", err)
	}
	defer file.Close()
	r := csv.NewReader(file)
	r.Comma = '\t'
	r.Read()
	for i := 0; ; i++ {
		_, err := ReadLine(r)
		if err != nil && err != io.EOF {
			continue
		} else if err == io.EOF {
			fmt.Println(i)
			break
		}
	}
	es, _ := elasticsearch.NewDefaultClient()
	log.Println(elasticsearch.Version)
	log.Println(es.Info())
}
