package user_domain

import "context"

type VerifyEmailUniquenessService struct {
	userRepo IUserRepository
}

func NewVerifyEmailUniquenessService(userRepository IUserRepository) *VerifyEmailUniquenessService {
	return &VerifyEmailUniquenessService{
		userRepo: userRepository,
	}
}

func (s *VerifyEmailUniquenessService) Run(ctx context.Context, email string) error {
	exists, err := s.userRepo.ExistsByEmail(ctx, email)
	if err != nil {
		return err
	}
	if exists {
		return ErrUserEmailAlreadyExists
	}
	return nil
}
