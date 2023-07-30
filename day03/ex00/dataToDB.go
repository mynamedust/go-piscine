package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"github.com/elastic/go-elasticsearch/v8"
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
	cfg := elasticsearch.Config{
        Addresses: []string{
            "http://localhost:9200",
        },
	}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
  		fmt.Printf("Error creating the client: %s", err)
	}
	es.Get("places","1")
}
