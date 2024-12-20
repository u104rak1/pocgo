package authentication

import (
	"context"

	userDomain "github.com/u104rak1/pocgo/internal/domain/user"
	idVO "github.com/u104rak1/pocgo/internal/domain/value_object/id"
)

type IAuthenticationService interface {
	VerifyUniqueness(ctx context.Context, userID idVO.UserID) error
	Authenticate(ctx context.Context, email, password string) (*idVO.UserID, error)
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

func (s *authenticationService) VerifyUniqueness(ctx context.Context, userID idVO.UserID) error {
	exists, err := s.authRepo.ExistsByUserID(ctx, userID)
	if err != nil {
		return err
	}
	if exists {
		return ErrAlreadyExists
	}
	return nil
}

func (s *authenticationService) Authenticate(ctx context.Context, email, password string) (*idVO.UserID, error) {
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
