package authentication

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/ucho456job/pocgo/internal/config"
)

type IAuthenticationService interface {
	VerifyUniqueness(ctx context.Context, userID string) error
	GenerateAccessToken(ctx context.Context, userID string) (string, error)
	GetUserIDFromAccessToken(ctx context.Context, accessToken string) (string, error)
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

func (s *authenticationService) GenerateAccessToken(ctx context.Context, userID string) (string, error) {
	env := config.NewEnv()
	jwtSecret := []byte(env.JWT_SECRET_KEY)

	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func (s *authenticationService) GetUserIDFromAccessToken(ctx context.Context, accessToken string) (string, error) {
	env := config.NewEnv()
	jwtSecret := []byte(env.JWT_SECRET_KEY)

	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrUnexpectedSigningMethod
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

	return "", ErrAuthenticationFailed
}
