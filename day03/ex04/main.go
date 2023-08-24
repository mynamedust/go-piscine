package main

import (
	"day03/ex04/db"
	"fmt"
	"net/http"
	"time"
)

const (
	userId     = "id328228"
	signingKey = "khimajKmsdf34#s48"
	tokenTTL   = 12 * time.Hour
)

func main() {
	client, err := db.CreateClient()
	if err != nil {
		return
	}
	esStore := db.NewElasticStore(client)

	jwtMiddleware := db.CreateJwtMiddleWare(signingKey)

	http.Handle("/api/recommend", jwtMiddleware.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
		responseData := db.CreatePageData(closestPlaces)
		db.CreatePage(*responseData, w)
	})))

	http.HandleFunc("/api/get_token", func(w http.ResponseWriter, r *http.Request) {
		if err := db.CheckJwtRequest(w, r); err != nil {
			return
		}
		token, err := db.GenerateJwtToken(w, r, userId, signingKey, tokenTTL)
		if err != nil {
			fmt.Println("Token creation failed: ", err)
			return
		}
		responseData := db.CreateJwtPageData(token)
		db.CreateJwtPage(*responseData, w)
	})

	http.ListenAndServe(":8888", nil)
	return
}
