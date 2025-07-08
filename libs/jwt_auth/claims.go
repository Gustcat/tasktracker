package jwt_auth

import "github.com/dgrijalva/jwt-go"

type UserClaims struct {
	jwt.StandardClaims
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Role int32  `json:"role"`
}
