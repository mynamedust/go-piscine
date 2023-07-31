package main

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esutil"
	"io"
	"os"
	"strconv"
	"strings"
	"sync/atomic"
)

func main() {
	var success uint64
	indexName := "places"
	mappingFileName := "../materials/schema.json"
	es, err := createClient()
	if err != nil {
		return
	}
	indexStatus := indexCreate(es, mappingFileName, indexName)
	if indexStatus != true {
		return
	}
	bulker, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Index:  indexName,
		Client: es,
	})
	actions, err := createActionList(indexName, &success)
	if err != nil {
		return
	}
	bulkErr := bulkActionAdd(actions, bulker)
	if bulkErr != nil {
		return
	}
	fmt.Printf("requests completed: %d/%d success.\n", len(actions), int(success))
}

func createClient() (*elasticsearch.Client, error) {
	cfg := elasticsearch.Config{
		Addresses: []string{
			"http://localhost:9200",
		},
	}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		fmt.Printf("error: client creating failed: %s\n", err)
		return nil, err
	}
	_, err = es.Info()
	if err != nil {
		fmt.Println("error: client.Info(); connection failed: ", err)
		return nil, err
	}
	fmt.Printf("elasticsearch: client connected %s\n", cfg.Addresses[0])
	return es, err
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
		return true
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

func createActionList(indexName string, success *uint64) (actions []esutil.BulkIndexerItem, err error) {
	file, err := os.Open("../materials/data.csv")
	if err != nil {
		fmt.Println("error: file opening failed: ", err)
		return
	}
	defer file.Close()
	r := csv.NewReader(file)
	r.Comma = '\t'
	r.Read()
	for i := 0; ; i++ {
		data, err := ReadLine(r)
		if err != nil && err != io.EOF {
			continue
		} else if err == io.EOF {
			break
		}
		body, err := json.Marshal(data)
		if err != nil {
			fmt.Println("error: data marshalling failed:", err)
			continue
		}
		action := createBulkItem(data, body, success)
		actions = append(actions, action)
	}
	return
}

func createBulkItem(data *RestaurantsType, body []byte, success *uint64) esutil.BulkIndexerItem {
	return esutil.BulkIndexerItem{
		Action:     "index",
		DocumentID: strconv.Itoa(data.Id),
		Body:       bytes.NewReader(body),
		OnSuccess: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem) {
			atomic.AddUint64(success, 1)
		},
		OnFailure: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem, err error) {
			if err != nil {
				fmt.Printf("ERROR: %s\n", err)
			} else {
				fmt.Printf("ERROR: %s: %s\n", res.Error.Type, res.Error.Reason)
			}
		},
	}
}

func bulkActionAdd(actions []esutil.BulkIndexerItem, bulker esutil.BulkIndexer) error {
	var err error
	for _, action := range actions {
		err = bulker.Add(
			context.Background(),
			action,
		)
	}
	if err != nil {
		fmt.Printf("Unexpected error: %s\n", err)
		return err
	}
	if err := bulker.Close(context.Background()); err != nil {
		fmt.Printf("Unexpected error: %s\n", err)
		return err
	}
	return nil
}
