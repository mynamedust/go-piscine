package main

import (
	"day03/ex02/db"
	"net/http"
)

func main() {
	client, err := db.CreateClient()
	if err != nil {
		return
	}
	esStore := db.NewElasticStore(client)
	http.HandleFunc("/api/places", func(w http.ResponseWriter, r *http.Request) {
		limit := 5
		pageNumber := db.CheckPageRequest(w, r)
		if pageNumber == -1 {
			return
		}
		places, count, err := esStore.GetPlaces((pageNumber-1)*limit, limit)
		if err != nil {
			return
		}
		newPage := db.CreatePageData(count, places, pageNumber, limit, w)
		db.CreatePage(*newPage, w, r)
	})
	http.ListenAndServe(":8888", nil)
	return
}
