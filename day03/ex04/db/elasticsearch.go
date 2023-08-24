package db

import (
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
)

type ElasticStore struct {
	esClient *elasticsearch.Client
}

func CreateClient() (*elasticsearch.Client, error) {
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
