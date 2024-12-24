package inmemory

import (
	"context"
	"sync"

	accountDomain "github.com/u104rak1/pocgo/internal/domain/account"
	idVO "github.com/u104rak1/pocgo/internal/domain/value_object/id"
)

type accountInMemoryRepository struct {
	mu       sync.RWMutex
	accounts map[string]*accountDomain.Account
}

func NewAccountInMemoryRepository() accountDomain.IAccountRepository {
	return &accountInMemoryRepository{
		accounts: make(map[string]*accountDomain.Account),
	}
}

func (r *accountInMemoryRepository) Save(ctx context.Context, account *accountDomain.Account) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.accounts[account.IDString()] = account
	return nil
}

func (r *accountInMemoryRepository) FindByID(ctx context.Context, id idVO.AccountID) (*accountDomain.Account, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	account, exists := r.accounts[id.String()]
	if !exists {
		return nil, nil
	}
	return account, nil
}

func (r *accountInMemoryRepository) CountByUserID(ctx context.Context, userID idVO.UserID) (int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	count := 0
	for _, account := range r.accounts {
		if account.UserIDString() == userID.String() {
			count++
		}
	}
	return count, nil
}
