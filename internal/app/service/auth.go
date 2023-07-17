package service

import (
	"github.com/MrTomSawyer/loyalty-system/internal/app/logger"
	"github.com/MrTomSawyer/loyalty-system/internal/app/models"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type AuthService struct {
	secret string
	exp    int
}

func NewAuthService(secret string, exp int) *AuthService {
	return &AuthService{
		secret: secret,
		exp:    exp,
	}
}

func (a AuthService) JWT(id int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, models.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(a.exp))),
		},
		UserID: id,
	})

	tokenString, err := token.SignedString([]byte(a.secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (a AuthService) GetUserID(token string) (int, error) {
	claims := models.Claims{}

	t, err := jwt.ParseWithClaims(token, &claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(a.secret), nil
	})

	if err != nil {
		logger.Log.Errorf("failed to parse token: %v", err)
		return 0, err
	}

	if !t.Valid {
		logger.Log.Errorf("token is not valid: %v", err)
		return 0, err
	}

	return claims.UserID, nil
}
