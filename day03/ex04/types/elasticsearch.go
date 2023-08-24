package types

type Response struct {
	Aggregations struct {
		TotalCount struct {
			Value float64 `json:"value"`
		} `json:"total_count"`
	} `json:"aggregations"`
	Hits struct {
		Hits []struct {
			Source Place `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}
