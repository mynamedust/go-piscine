package types

import "github.com/dgrijalva/jwt-go"

type JwtPage struct {
	Token string `json:"token"`
}

type JwtTokenClaims struct {
	jwt.StandardClaims
	UserId string `json:"user_id"`
}
