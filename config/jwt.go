package config

import (
	"github.com/golang-jwt/jwt/v5"
)

type TokenClaims[T any] struct {
	Payload T
	jwt.RegisteredClaims
}

func (c TokenClaims[T]) GenerateToken(key []byte) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return token.SignedString(key)
}

func VerifyToken[T any](tokenString string, key []byte, claimData T) (T, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims[T]{}, func(t *jwt.Token) (interface{}, error) {
		return key, nil
	})

	if claims, ok := token.Claims.(*TokenClaims[T]); ok && token.Valid {
		return claims.Payload, nil
	}

	return TokenClaims[T]{}.Payload, err
}

