package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"

	//"strings"
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
	fmt.Println(es)
	mapping, status := mappingCreate("../materials/schema.json")
	if status == 1 {
		return
	}
	res, err := es.Indices.Create("places", es.Indices.Create.WithBody(strings.NewReader(mapping)))
	fmt.Println(res)
}

func mappingCreate(fileName string) (string, int) {
	file, err := os.Open(fileName)
	defer file.Close()
	if err != nil {
		fmt.Println("error: file opening failed: ", err)
		return "", 1
	}
	data, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("error: file reading failed: ", err)
		return "", 1
	}
	mapping := fmt.Sprintf("{\n\"mappings\": %s}", string(data))
	return mapping, 0
}
