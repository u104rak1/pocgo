package user

import (
	"context"

	idVO "github.com/u104rak1/pocgo/internal/domain/value_object/id"
)

type IUserService interface {
	VerifyEmailUniqueness(ctx context.Context, email string) error
	EnsureUserExists(ctx context.Context, id idVO.UserID) error
	FindUser(ctx context.Context, id idVO.UserID) (*User, error)
}

type userService struct {
	userRepo IUserRepository
}

func NewService(userRepository IUserRepository) IUserService {
	return &userService{
		userRepo: userRepository,
	}
}

func (s *userService) VerifyEmailUniqueness(ctx context.Context, email string) error {
	exists, err := s.userRepo.ExistsByEmail(ctx, email)
	if err != nil {
		return err
	}
	if exists {
		return ErrEmailAlreadyExists
	}
	return nil
}

func (s *userService) EnsureUserExists(ctx context.Context, id idVO.UserID) error {
	exists, err := s.userRepo.ExistsByID(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return ErrNotFound
	}
	return nil
}

func (s *userService) FindUser(ctx context.Context, id idVO.UserID) (*User, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrNotFound
	}
	return user, nil
}
