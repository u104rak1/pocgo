package authentication

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/ucho456job/pocgo/internal/config"
)

type AccessTokenService struct{}

func NewAccessTokenService() *AccessTokenService {
	return &AccessTokenService{}
}

func (s *AccessTokenService) Generate(ctx context.Context, userID string) (string, error) {
	env := config.NewEnv()
	jwtSecret := []byte(env.JWT_SECRET_KEY)

	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func (s *AccessTokenService) GetUserID(ctx context.Context, tokenString string) (string, error) {
	env := config.NewEnv()
	jwtSecret := []byte(env.JWT_SECRET_KEY)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtSecret, nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if userID, ok := claims["userID"].(string); ok {
			return userID, nil
		}
	}

	return "", errors.New("invalid token or missing userID")
}
