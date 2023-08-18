package db

import (
	"day03/ex01/types"
	"html/template"
	"net/http"
	"strconv"
)

func CreatePage(htmlName string, pageData types.HtmlPage, w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(htmlName)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err = tmpl.Execute(w, pageData)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
	return
}

func CheckPageRequest(w http.ResponseWriter, r *http.Request) int {
	pageStr := r.URL.Query().Get("page")
	if pageStr == "" {
		http.Error(w, "Missing 'page' parameter", http.StatusBadRequest)
		return -1
	}
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		http.Error(w, "Invalid 'page' value", http.StatusBadRequest)
		return -1
	}
	return page
}
