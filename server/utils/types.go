package utils

import "github.com/golang-jwt/jwt/v5"

type AuthorizedUserInfo struct {
	ID			string		`json:"id"`
	Username	string		`json:"username"`
	jwt.RegisteredClaims
}
