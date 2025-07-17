package utils

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// GenerateToken creates a JWT token with username, user_id, and role
func GenerateToken(username string, userID int, role string) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", errors.New("JWT_SECRET is not set in environment variables")
	}

	claims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"role":     role,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// VerifyToken parses and validates JWT, returning ID, username, and role
func VerifyToken(tokenString string) (int, string, string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return 0, "", "", errors.New("JWT_SECRET is not set")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Check signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})

	if err != nil || !token.Valid {
		return 0, "", "", errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, "", "", errors.New("invalid claims format")
	}

	// Extract and cast user ID
	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return 0, "", "", errors.New("user_id not found or invalid")
	}
	userID := int(userIDFloat)

	username, _ := claims["username"].(string)
	role, _ := claims["role"].(string)

	return userID, username, role, nil
}
