package authentication

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	userDomain "github.com/ucho456job/pocgo/internal/domain/user"
	"github.com/ucho456job/pocgo/pkg/timer"
)

type IAuthenticationService interface {
	VerifyUniqueness(ctx context.Context, userID string) error
	GenerateAccessToken(ctx context.Context, userID string, jwtSecretKey []byte) (string, error)
	GetUserIDFromAccessToken(ctx context.Context, accessToken string, jwtSecretKey []byte) (string, error)
	Authenticate(ctx context.Context, email, password string) (userID string, err error)
}

type authenticationService struct {
	authRepo IAuthenticationRepository
	userRepo userDomain.IUserRepository
}

func NewService(authenticationRepository IAuthenticationRepository, userRepository userDomain.IUserRepository) IAuthenticationService {
	return &authenticationService{
		authRepo: authenticationRepository,
		userRepo: userRepository,
	}
}

func (s *authenticationService) VerifyUniqueness(ctx context.Context, userID string) error {
	exists, err := s.authRepo.ExistsByUserID(ctx, userID)
	if err != nil {
		return err
	}
	if exists {
		return ErrAlreadyExists
	}
	return nil
}

func (s *authenticationService) GenerateAccessToken(ctx context.Context, userID string, jwtSecretKey []byte) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": timer.Now().Add(time.Hour * 24).Unix(),
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

	return "", ErrInvalidAccessToken
}

func (s *authenticationService) Authenticate(ctx context.Context, email, password string) (userID string, err error) {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		if err == userDomain.ErrNotFound {
			return "", ErrAuthenticationFailed
		}
		return "", err
	}

	auth, err := s.authRepo.FindByUserID(ctx, user.ID())
	if err != nil {
		if err == ErrNotFound {
			return "", ErrAuthenticationFailed
		}
		return "", err
	}

	if err := auth.ComparePassword(password); err != nil {
		return "", ErrAuthenticationFailed
	}

	return user.ID(), nil
}
