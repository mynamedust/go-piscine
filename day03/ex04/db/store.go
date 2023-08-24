package db

import (
	"day03/ex04/types"
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v8"
	"strconv"
	"strings"
)

type Store interface {
	GetClosestPlaces(limit int) ([]types.Place, int, error)
}

func NewElasticStore(esClient *elasticsearch.Client) *ElasticStore {
	return &ElasticStore{esClient: esClient}
}

func (es *ElasticStore) GetClosestPlaces(limit int, lat float64, lon float64) (places []types.Place, err error) {
	var response types.Response
	res, err := es.esClient.Search(es.esClient.Search.WithBody(strings.NewReader(`{
		"from": 0,
		"size": `+strconv.Itoa(limit)+`,
		"sort": [
        {
            "_geo_distance": {
                "location": {
                    "lat": `+strconv.FormatFloat(lat, 'f', 3, 64)+`,
                    "lon": `+strconv.FormatFloat(lon, 'f', 3, 64)+`
                },
                "order": "asc",
                "unit": "km",
                "mode": "min",
                "distance_type": "arc",
                "ignore_unmapped": true
            }
        }
    ]
	}`)),
		es.esClient.Search.WithPretty(),
	)
	defer res.Body.Close()
	if err != nil {
		return nil, err
	}
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}
	for _, hit := range response.Hits.Hits {
		places = append(places, hit.Source)
	}
	return
}
