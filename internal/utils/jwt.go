package utils

import (
	"errors"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/lucaserm/ecom/internal/env"
)

var ErrInvalidToken = errors.New("invalid token")

func jwtSecret() string {
	return env.GetString("JWT_SECRET", "somerandomsecret")
}

func CreateJWT(userId string) (string, error) {
	secret := jwtSecret()
	expInSeconds := env.GetInt("JWT_EXP", 60*60*24*7) // a week
	exp := time.Second * time.Duration(expInSeconds)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userId,
		"exp":    time.Now().Add(exp).Unix(),
	})

	tokenStr, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

// VerifyJWT parses and validates an HS256 token, returning the userId claim.
func VerifyJWT(tokenStr string) (string, error) {
	secret := jwtSecret()

	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(secret), nil
	})
	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", ErrInvalidToken
	}

	userId, ok := claims["userId"].(string)
	if !ok || userId == "" {
		return "", ErrInvalidToken
	}

	return userId, nil
}
