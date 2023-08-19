package main

import (
	"day03/ex01/db"
	"day03/ex01/types"
	"net/http"
)

func main() {
	htmlName := "./templates/index.html"
	client, err := db.CreateClient()
	if err != nil {
		return
	}
	esStore := db.NewElasticStore(client)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		limit := 5
		pageNumber := db.CheckPageRequest(w, r)
		if pageNumber == -1 {
			return
		}
		places, count, err := esStore.GetPlaces((pageNumber-1)*limit, limit)
		if err != nil {
			return
		}
		newPage := types.HtmlPage{
			TotalCount:   count,
			Places:       places,
			HasPrevious:  pageNumber > 1,
			HasNext:      pageNumber*limit < count,
			PreviousPage: pageNumber - 1,
			NextPage:     pageNumber + 1,
			LastPage:     count / limit,
		}
		if count%limit > 0 {
			newPage.LastPage += 1
		}
		if pageNumber > newPage.LastPage {
			http.Error(w, "Invalid 'page' value: 'foo'", 400)
			return
		}
		db.CreatePage(htmlName, newPage, w, r)
	})
	http.ListenAndServe(":8888", nil)
	return
}
