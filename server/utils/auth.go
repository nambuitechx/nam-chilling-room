package utils

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

func ValidateTokenString(tokenString string) (*AuthorizedUserInfo, error) {
	claims := &AuthorizedUserInfo{}

	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(t *jwt.Token) (any, error) {
			// check signing method
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return []byte("secret"), nil
		},
	)

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
