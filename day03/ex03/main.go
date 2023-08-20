package main

import (
	"day03/ex03/db"
	"fmt"
	"net/http"
)

func main() {
	client, err := db.CreateClient()
	if err != nil {
		return
	}
	esStore := db.NewElasticStore(client)
	http.HandleFunc("/api/recommend", func(w http.ResponseWriter, r *http.Request) {
		limit := 3
		lat, lon, err := db.CheckGeoSearchRequest(w, r)
		if err != nil {
			return
		}
		closestPlaces, err := esStore.GetClosestPlaces(limit, lat, lon)
		if err != nil {
			fmt.Println("Page data creation error.")
			return
		}
		pageData := db.CreatePageData(closestPlaces)
		db.CreatePage(*pageData, w, r)
	})
	http.ListenAndServe(":8888", nil)
	return
}
