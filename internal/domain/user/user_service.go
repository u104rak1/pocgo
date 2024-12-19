package user

import "context"

type IUserService interface {
	VerifyEmailUniqueness(ctx context.Context, email string) error
	EnsureUserExists(ctx context.Context, id UserID) error
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

func (s *userService) EnsureUserExists(ctx context.Context, id UserID) error {
	exists, err := s.userRepo.ExistsByID(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return ErrNotFound
	}
	return nil
}
