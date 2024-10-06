package user_domain

import "context"

type IsEmailDuplicateService struct {
	userRepo IUserRepository
}

func NewIsEmailDuplicateService(userRepository IUserRepository) *IsEmailDuplicateService {
	return &IsEmailDuplicateService{
		userRepo: userRepository,
	}
}

func (s *IsEmailDuplicateService) IsEmailDuplicate(ctx context.Context, email string) (bool, error) {
	exists, err := s.userRepo.ExistsByEmail(ctx, email)
	if err != nil {
		return false, ErrUserEmailAlreadyExists
	}
	return exists, nil
}
