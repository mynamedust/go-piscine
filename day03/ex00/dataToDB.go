package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

func main() {
	indexName := "places"
	mappingFileName := "../materials/schema.json"
	file, err := os.Open("../materials/data.csv")
	if err != nil {
		fmt.Println("error: file opening failed: ", err)
	}
	actions, err := createActionList(file, indexName)
	fmt.Println(len(actions))
	if err != nil {
		return
	}
	cfg := elasticsearch.Config{
		Addresses: []string{
			"http://localhost:9200",
		},
	}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		fmt.Printf("error: client creating failed: %s\n", err)
	}
	_, err = es.Info()
	if err != nil {
		fmt.Println("error: client.Info(); connection failed: ", err)
		return
	} else {
		fmt.Println("elasticsearch: client connected.")
	}
	indexStatus := indexCreate(es, mappingFileName, indexName)
	if indexStatus != true {
		return
	}
}

func indexExist(index string, client *elasticsearch.Client) (bool, error) {
	res, err := client.Indices.Exists([]string{index})
	if err != nil {
		return false, err
	}
	defer res.Body.Close()
	if res.IsError() {
		return false, nil
	}
	fmt.Println("error: index creating failed: index already exists.")
	return true, nil
}

func indexCreate(client *elasticsearch.Client, mapping, indexName string) bool {
	mapping, status := mappingCreate(mapping)
	if status == 1 {
		return false
	}
	exists, _ := indexExist(indexName, client)
	if exists {
		return false
	}
	res, _ := client.Indices.Create(indexName, client.Indices.Create.WithBody(strings.NewReader(mapping)))
	if res.StatusCode != 200 {
		fmt.Println("error: index creating failed.")
	}
	return true
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

func createActionList(file *os.File, indexName string) (actions []esapi.IndexRequest, err error) {
	var success, i int
	defer file.Close()
	r := csv.NewReader(file)
	r.Comma = '\t'
	r.Read()
	for i = 0; ; i++ {
		data, err := ReadLine(r)
		if err != nil && err != io.EOF {
			continue
		} else if err == io.EOF {
			break
		}
		body, err := json.Marshal(data)
		if err != nil {
			fmt.Println("error: data marshalling failed:" , err)
			continue
		}
		action := esapi.IndexRequest{
			Index: indexName,
			DocumentID: strconv.Itoa(data.Id),
			Body: bytes.NewReader(body),
		}
		actions = append(actions, action)
		success++
	}
	fmt.Printf("action list creation finished: %d/%d success.\n", i, success)
	return
}
