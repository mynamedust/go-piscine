package db

import (
	"day03/ex04/types"
	"encoding/json"
	"net/http"
)

func CreatePage(pageData types.JsonPage, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "	")
	if err := encoder.Encode(pageData); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	return
}

func CreatePageData(places []types.Place) *types.JsonPage {
	newPage := types.JsonPage{
		Name:   "Recommendation",
		Places: places,
	}
	return &newPage
}
