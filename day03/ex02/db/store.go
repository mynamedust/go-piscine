package db

import (
	"day03/ex01/types"
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v8"
	"strconv"
	"strings"
)

type Store interface {
	GetPlaces(limit int, offset int) ([]types.Place, int, error)
}

func NewElasticStore(esClient *elasticsearch.Client) *ElasticStore {
	return &ElasticStore{esClient: esClient}
}

func (es *ElasticStore) GetPlaces(limit int, offset int) (places []types.Place, count int, err error) {
	var response types.Response
	res, err := es.esClient.Search(es.esClient.Search.WithBody(strings.NewReader(`{
		"from": `+strconv.Itoa(limit)+`,
		"size": `+strconv.Itoa(offset)+`,
		"query": {
			"match": {
				"_index": "places"
			}
		},
		"aggs": {
			"total_count": {
				"value_count": {
					"field": "_id"
				}
			}
		}
	}`)),
		es.esClient.Search.WithPretty(),
	)
	defer res.Body.Close()
	if err != nil {
		return nil, 0, err
	}
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, 0, err
	}
	for _, hit := range response.Hits.Hits {
		places = append(places, hit.Source)
	}
	count = int(response.Aggregations.TotalCount.Value)
	return
}
