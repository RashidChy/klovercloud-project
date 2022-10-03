package dto

import (
	"github.com/golang-jwt/jwt"
)

// JwtCustomClaims are custom claims extending default ones.
type JwtCustomClaims struct {
	Name      string `json:"Name"`
	Email     string `json:"Email"`
	SlackHook string `json:"SlackHook"`
	jwt.StandardClaims
}

type Token struct {
	Token string `json:"token"`
}
