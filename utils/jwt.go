package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func GenerateJWT(email string, role string, userID, accountID string) (string, error) {

	claims := jwt.MapClaims{
		"email":     email,
		"role":      role,
		"userID":    userID,
		"accountID": accountID,                             // Placeholder for user ID, to be set when generating the token
		"exp":       time.Now().Add(time.Hour * 24).Unix(), // expire in 24 hours
		"iat":       time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(jwtSecret)
}

func VerifyJWT(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	// convert claims to MapClaims
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims, nil
	}

	return nil, err
}
