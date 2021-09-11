package main

import (
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// Data to be encrypted with JWT
type Claims struct {
	UserId string `json:"userId,omitempty"`
	jwt.StandardClaims
}

// Generates a JWT token containing userId
func GenerateToken(userId string, expHours int64) (string, error) {
	rawToken := jwt.NewWithClaims(jwt.SigningMethodHS512, Claims{
		UserId: userId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * time.Duration(expHours)).Unix(),
			IssuedAt:  time.Now().Unix(),
			Subject:   "auth_token",
		},
	})

	return rawToken.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

// Verifies JWT validity & returns data encrypted within
func VerifyToken(token string) (*jwt.Token, Claims, error) {
	var tokenData Claims

	parsedToken, err := jwt.ParseWithClaims(token, &tokenData, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	return parsedToken, tokenData, err
}
