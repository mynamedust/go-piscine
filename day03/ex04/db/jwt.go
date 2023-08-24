package db

import (
	"day03/ex04/types"
	"encoding/json"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"time"
)

func CheckJwtRequest(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodGet {
		err := errors.New("Method not allowed")
		http.Error(w, err.Error(), http.StatusMethodNotAllowed)
		return err
	}
	return nil
}

func GenerateJwtToken(w http.ResponseWriter, r *http.Request, userId, signingKey string, timeTTL time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &types.JwtTokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(timeTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		userId,
	})
	return token.SignedString([]byte(signingKey))
}

func CreateJwtPageData(token string) *types.JwtPage {
	newPage := types.JwtPage{
		Token: token,
	}
	return &newPage
}

func CreateJwtPage(data types.JwtPage, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "	")
	if err := encoder.Encode(data); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
	return
}
