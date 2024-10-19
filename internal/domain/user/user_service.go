package user

import "context"

type IUserService interface {
	VerifyEmailUniqueness(ctx context.Context, email string) error
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
		return ErrUserEmailAlreadyExists
	}
	return nil
}
