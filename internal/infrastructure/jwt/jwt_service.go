package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	authApp "github.com/u104rak1/pocgo/internal/application/authentication"
	"github.com/u104rak1/pocgo/pkg/timer"
)

var (
	ErrUnexpectedSigningMethod = errors.New("unexpected signing method")
	ErrInvalidAccessToken      = errors.New("invalid access token")
	ErrUserIDMissing           = errors.New("user id missing")
)

type jwtService struct {
	secretKey []byte
}

func NewService(secretKey []byte) authApp.IJWTService {
	return &jwtService{
		secretKey: secretKey,
	}
}

func (s *jwtService) GenerateAccessToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": timer.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secretKey)
}

func (s *jwtService) GetUserIDFromAccessToken(accessToken string) (string, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrUnexpectedSigningMethod
		}
		return s.secretKey, nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if userID, ok := claims["sub"].(string); ok {
			return userID, nil
		}
	}

	return "", ErrInvalidAccessToken
}
