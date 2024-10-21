package authentication

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type IAuthenticationService interface {
	VerifyUniqueness(ctx context.Context, userID string) error
	GenerateAccessToken(ctx context.Context, userID string, jwtSecretKey []byte) (string, error)
	GetUserIDFromAccessToken(ctx context.Context, accessToken string, jwtSecretKey []byte) (string, error)
}

type authenticationService struct {
	authenticationRepo IAuthenticationRepository
}

func NewService(authenticationRepository IAuthenticationRepository) IAuthenticationService {
	return &authenticationService{
		authenticationRepo: authenticationRepository,
	}
}

func (s *authenticationService) VerifyUniqueness(ctx context.Context, userID string) error {
	exists, err := s.authenticationRepo.ExistsByUserID(ctx, userID)
	if err != nil {
		return err
	}
	if exists {
		return ErrAuthenticationAlreadyExists
	}
	return nil
}

func (s *authenticationService) GenerateAccessToken(ctx context.Context, userID string, jwtSecretKey []byte) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecretKey)
}

func (s *authenticationService) GetUserIDFromAccessToken(ctx context.Context, accessToken string, jwtSecretKey []byte) (string, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrUnexpectedSigningMethod
		}
		return jwtSecretKey, nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if userID, ok := claims["sub"].(string); ok {
			return userID, nil
		}
	}

	return "", ErrAuthenticationFailed
}
