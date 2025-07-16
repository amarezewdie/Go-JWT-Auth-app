package utils

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// GenerateToken creates a JWT token with username and user_id
func GenerateToken(username string, id int) (string, error) {

	// Payload (claims) for the JWT token
	claims := jwt.MapClaims{
		"username": username,
		"user_id":  1,
		"exp":      time.Now().Add(time.Hour * 1).Unix(),
	}

	// Create the token with claims and sign it
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", errors.New("JWT_SECRET is not set")
	}

	return token.SignedString([]byte(secret))

}

func VerifyToken(tokenString string) (int, string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return 0, "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID := int(claims["user_id"].(float64))
		userName := claims["username"].(string)

		return userID, userName, nil
	}

	return 0, "", err
}
