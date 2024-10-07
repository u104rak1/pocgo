package authentication_domain

import "context"

type VerifyAuthenticationUniquenessService struct {
	authenticationRepo IAuthenticationRepository
}

func NewVerifyEmailUniquenessService(authenticationRepository IAuthenticationRepository) *VerifyAuthenticationUniquenessService {
	return &VerifyAuthenticationUniquenessService{
		authenticationRepo: authenticationRepository,
	}
}

func (s *VerifyAuthenticationUniquenessService) Run(ctx context.Context, userID string) error {
	a, err := s.authenticationRepo.FindByUserID(ctx, userID)
	if err != nil {
		return err
	}
	if a != nil {
		return ErrAuthenticationAlreadyExists
	}
	return nil
}
