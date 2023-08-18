package main

import (
	"day03/ex01/db"
	"day03/ex01/types"
	"fmt"
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
		fmt.Println(r.URL, ": page change handled")
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
		db.CreatePage(htmlName, newPage, w, r)
	})
	fmt.Println("i am in the end")
	http.ListenAndServe(":8888", nil)
	return
}
