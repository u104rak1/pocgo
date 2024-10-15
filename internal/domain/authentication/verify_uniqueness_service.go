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
	exists, err := s.authenticationRepo.ExistsByUserID(ctx, userID)
	if err != nil {
		return err
	}
	if exists {
		return ErrAuthenticationAlreadyExists
	}
	return nil
}
