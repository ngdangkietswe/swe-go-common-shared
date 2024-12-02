package util

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
)

// ParseToken is a function that parses a JWT token.
func ParseToken(jwtToken, jwtSecret string) (*jwt.MapClaims, error) {
	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !token.Valid || !ok {
		return nil, err
	}

	return &claims, nil
}
