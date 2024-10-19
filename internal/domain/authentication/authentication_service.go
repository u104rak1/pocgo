package authentication

import "context"

type IAuthenticationService interface {
	VerifyUniqueness(ctx context.Context, userID string) error
}

type authenticationService struct {
	authenticationRepo IAuthenticationRepository
}

func NewService(authenticationRepository IAuthenticationRepository) IAuthenticationService {
	return &authenticationService{
		authenticationRepo: authenticationRepository,
	}
}

func (s *authenticationService) VerifyUniqueness(ctx context.Context, userID string) error {
	exists, err := s.authenticationRepo.ExistsByUserID(ctx, userID)
	if err != nil {
		return err
	}
	if exists {
		return ErrAuthenticationAlreadyExists
	}
	return nil
}
