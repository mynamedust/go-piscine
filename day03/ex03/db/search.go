package db

import (
	"errors"
	"net/http"
	"strconv"
)

func CheckGeoSearchRequest(w http.ResponseWriter, r *http.Request) (lat float64, lon float64, err error) {
	strKeys := r.URL.Query()
	latStr, lonStr := strKeys.Get("lat"), strKeys.Get("lon")
	if latStr == "" || lonStr == "" {
		err = errors.New("Missing 'lat' and 'lon' parameters")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	lat, err = strconv.ParseFloat(latStr, 64)
	if err != nil {
		http.Error(w, "Invalid 'lat' value", http.StatusBadRequest)
		return
	}
	lon, err = strconv.ParseFloat(lonStr, 64)
	if err != nil {
		http.Error(w, "Invalid 'lon' value", http.StatusBadRequest)
	}
	if r.Method != http.MethodGet {
		err = errors.New("Method not allowed")
		http.Error(w, err.Error(), http.StatusMethodNotAllowed)
		return
	}
	return
}
