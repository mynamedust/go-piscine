package db

import (
	"day03/ex02/types"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func CreatePage(pageData types.JsonPage, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "	")
	if err := encoder.Encode(pageData); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	return
}

func CreatePageData(count int, places []types.Place, pageNumber int, limit int, w http.ResponseWriter) *types.JsonPage {
	newPage := types.JsonPage{
		Name:         "Places",
		TotalCount:   count,
		Places:       places,
		PreviousPage: pageNumber - 1,
		NextPage:     pageNumber + 1,
		LastPage:     count / limit,
	}
	if count%limit > 0 {
		newPage.LastPage += 1
	}
	if newPage.NextPage > newPage.LastPage {
		newPage.NextPage = 0
	}
	if pageNumber > newPage.LastPage || pageNumber <= 0 {
		fmt.Println("page number ", pageNumber)
		w.WriteHeader(http.StatusBadRequest)
		newPage = types.JsonPage{
			Error: "Invalid 'page' value: 'foo'",
		}
	}
	return &newPage
}

func CheckPageRequest(w http.ResponseWriter, r *http.Request) int {
	pageStr := r.URL.Query().Get("page")
	if pageStr == "" {
		http.Error(w, "Missing 'page' parameter", http.StatusBadRequest)
		return -1
	}
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		http.Error(w, "Invalid 'page' value: 'foo'", 400)
		return -1
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return -1
	}
	fmt.Println(r.Host, ": page change handled")
	return page
}
