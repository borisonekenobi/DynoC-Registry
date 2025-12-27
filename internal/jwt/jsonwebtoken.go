package jwt

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

var jwtSecret string

func Init() {
	jwtSecret = os.Getenv("JWT_SECRET")
}

func CreateToken(userId pgtype.UUID) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userId,
		"exp":     time.Now().Add(time.Hour * 1).Unix(),
	})

	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func verifyToken(tokenString string) (jwt.MapClaims, error) {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

func GetTokenClaims(tokenString string) (pgtype.UUID, error) {
	if tokenString == "" {
		return pgtype.UUID{}, fmt.Errorf("no token")
	}
	tokenString = tokenString[len("Bearer "):]

	claims, err := verifyToken(tokenString)
	if err != nil {
		return pgtype.UUID{}, err
	}

	var userId pgtype.UUID
	err = userId.Scan(claims["user_id"])
	if err != nil {
		return pgtype.UUID{}, err
	}

	return userId, nil
}
