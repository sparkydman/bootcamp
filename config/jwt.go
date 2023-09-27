package config

import "github.com/golang-jwt/jwt/v5"

type TokenClaims[T any] struct {
	Payload T
	jwt.RegisteredClaims
}

func (c TokenClaims[T]) GenerateToken(key []byte)(string, error){
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return token.SignedString(key)
}