package account

import "context"

type IAccountService interface {
	// Checks whether the user has reached the maximum number of accounts.
	CheckLimit(ctx context.Context, userID string) error

	// Get the user's account and check the user ID and password. Password confirmation is optional and can be nil to skip password confirmation.
	GetAndAuthorize(ctx context.Context, accountID string, userID, password *string) (*Account, error)
}

type accountService struct {
	accountRepo IAccountRepository
}

func NewService(accountRepository IAccountRepository) IAccountService {
	return &accountService{
		accountRepo: accountRepository,
	}
}

func (s *accountService) CheckLimit(ctx context.Context, userID string) error {
	count, err := s.accountRepo.CountByUserID(ctx, userID)
	if err != nil {
		return err
	}
	if count >= MaxAccountLimit {
		return ErrLimitReached
	}

	return nil
}

func (s *accountService) GetAndAuthorize(ctx context.Context, accountID string, userID, password *string) (*Account, error) {
	account, err := s.accountRepo.FindByID(ctx, accountID)
	if err != nil {
		return nil, err
	}
	if account == nil {
		return nil, ErrNotFound
	}
	if userID != nil && account.UserID() != *userID {
		return nil, ErrUnauthorized
	}
	if password != nil {
		if err := account.ComparePassword(*password); err != nil {
			return nil, err
		}
	}

	return account, nil
}
