package user

import "context"

type IVerifyEmailUniquenessService interface {
	Run(ctx context.Context, email string) error
}

type verifyEmailUniquenessService struct {
	userRepo IUserRepository
}

func NewVerifyEmailUniquenessService(userRepository IUserRepository) IVerifyEmailUniquenessService {
	return &verifyEmailUniquenessService{
		userRepo: userRepository,
	}
}

func (s *verifyEmailUniquenessService) Run(ctx context.Context, email string) error {
	exists, err := s.userRepo.ExistsByEmail(ctx, email)
	if err != nil {
		return err
	}
	if exists {
		return ErrUserEmailAlreadyExists
	}
	return nil
}
