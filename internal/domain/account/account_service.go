package account

import (
	"context"

	idVO "github.com/u104rak1/pocgo/internal/domain/value_object/id"
)

type IAccountService interface {
	// ユーザーの口座数が上限に達しているかをチェックします。
	CheckLimit(ctx context.Context, userID idVO.UserID) error

	// ユーザーの口座を取得する。ユーザーIDとパスワードの確認はオプションであり、必要ない場合はnilを渡す。
	GetAndAuthorize(ctx context.Context, accountID idVO.AccountID, userID *idVO.UserID, password *string) (*Account, error)
}

type accountService struct {
	accountRepo IAccountRepository
}

func NewService(accountRepository IAccountRepository) IAccountService {
	return &accountService{
		accountRepo: accountRepository,
	}
}

func (s *accountService) CheckLimit(ctx context.Context, userID idVO.UserID) error {
	count, err := s.accountRepo.CountByUserID(ctx, userID)
	if err != nil {
		return err
	}
	if count >= MaxAccountLimit {
		return ErrLimitReached
	}

	return nil
}

func (s *accountService) GetAndAuthorize(ctx context.Context, accountID idVO.AccountID, userID *idVO.UserID, password *string) (*Account, error) {
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
