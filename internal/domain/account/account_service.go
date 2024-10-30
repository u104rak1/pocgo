package account

import "context"

type IAccountService interface {
	CheckLimit(ctx context.Context, userID string) error
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
