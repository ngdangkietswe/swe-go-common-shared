package util

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/ngdangkietswe/swe-go-common-shared/config"
	"github.com/ngdangkietswe/swe-go-common-shared/grpc/domain"
	"time"
)

type JwtUserClaims struct {
	UserId   string `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// GenerateToken is a function that generates a JWT token.
func GenerateToken(grpcUser *domain.GrpcUser, isRefresh bool) (string, error) {
	var tokenExp time.Duration
	if isRefresh {
		tokenExp = time.Second * time.Duration(config.GetInt("REFRESH_TOKEN_EXPIRATION", 7200))
	} else {
		tokenExp = time.Second * time.Duration(config.GetInt("JWT_EXPIRATION", 3600))
	}

	exp := time.Now().Add(tokenExp).Unix()

	mapClaims := make(jwt.MapClaims)
	mapClaims["sub"] = grpcUser.Id
	mapClaims["user"] = JwtUserClaims{
		UserId:   grpcUser.Id,
		Username: grpcUser.Username,
		Email:    grpcUser.Email,
	}
	mapClaims["iat"] = time.Now().Unix()
	mapClaims["nbf"] = time.Now().Unix()
	mapClaims["exp"] = exp

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, mapClaims).SignedString([]byte(config.GetString("JWT_SECRET", "")))

	if err != nil {
		return "", err
	}

	return token, nil
}

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
