package authentication

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	userDomain "github.com/u104rak1/pocgo/internal/domain/user"
	"github.com/u104rak1/pocgo/pkg/timer"
)

type IAuthenticationService interface {
	VerifyUniqueness(ctx context.Context, userID userDomain.UserID) error
	GenerateAccessToken(ctx context.Context, userID userDomain.UserID, jwtSecretKey []byte) (string, error)
	GetUserIDFromAccessToken(ctx context.Context, accessToken string, jwtSecretKey []byte) (*userDomain.UserID, error)
	Authenticate(ctx context.Context, email, password string) (*userDomain.UserID, error)
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

func (s *authenticationService) VerifyUniqueness(ctx context.Context, userID userDomain.UserID) error {
	exists, err := s.authRepo.ExistsByUserID(ctx, userID)
	if err != nil {
		return err
	}
	if exists {
		return ErrAlreadyExists
	}
	return nil
}

func (s *authenticationService) GenerateAccessToken(ctx context.Context, userID userDomain.UserID, jwtSecretKey []byte) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": timer.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecretKey)
}

func (s *authenticationService) GetUserIDFromAccessToken(ctx context.Context, accessToken string, jwtSecretKey []byte) (*userDomain.UserID, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrUnexpectedSigningMethod
		}
		return jwtSecretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if userID, ok := claims["sub"].(string); ok {
			userID := userDomain.UserID(userID)
			return &userID, nil
		}
	}

	return nil, ErrInvalidAccessToken
}

func (s *authenticationService) Authenticate(ctx context.Context, email, password string) (*userDomain.UserID, error) {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrAuthenticationFailed
	}

	auth, err := s.authRepo.FindByUserID(ctx, user.ID())
	if err != nil {
		return nil, err
	}
	if auth == nil {
		return nil, ErrAuthenticationFailed
	}

	if err := auth.ComparePassword(password); err != nil {
		return nil, ErrAuthenticationFailed
	}

	userID := user.ID()
	return &userID, nil
}
