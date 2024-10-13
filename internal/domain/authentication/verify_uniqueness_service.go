package authentication

import "context"

type IVerifyAuthenticationUniquenessService interface {
	Run(ctx context.Context, userID string) error
}

type verifyAuthenticationUniquenessService struct {
	authenticationRepo IAuthenticationRepository
}

func NewVerifyAuthenticationUniquenessService(authenticationRepository IAuthenticationRepository) IVerifyAuthenticationUniquenessService {
	return &verifyAuthenticationUniquenessService{
		authenticationRepo: authenticationRepository,
	}
}

func (s *verifyAuthenticationUniquenessService) Run(ctx context.Context, userID string) error {
	a, err := s.authenticationRepo.FindByUserID(ctx, userID)
	if err != nil {
		return err
	}
	if a != nil {
		return ErrAuthenticationAlreadyExists
	}
	return nil
}
