package inmemory

import (
	"context"
	"sync"

	authDomain "github.com/u104rak1/pocgo/internal/domain/authentication"
	idVO "github.com/u104rak1/pocgo/internal/domain/value_object/id"
)

type authenticationInMemoryRepository struct {
	mu              sync.RWMutex
	authentications map[string]*authDomain.Authentication
}

func NewAuthenticationInMemoryRepository() authDomain.IAuthenticationRepository {
	return &authenticationInMemoryRepository{
		authentications: make(map[string]*authDomain.Authentication),
	}
}

func (r *authenticationInMemoryRepository) Save(ctx context.Context, authentication *authDomain.Authentication) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.authentications[authentication.UserIDString()] = authentication
	return nil
}

func (r *authenticationInMemoryRepository) FindByUserID(ctx context.Context, userID idVO.UserID) (*authDomain.Authentication, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	auth, exists := r.authentications[userID.String()]
	if !exists {
		return nil, nil
	}
	return auth, nil
}

func (r *authenticationInMemoryRepository) ExistsByUserID(ctx context.Context, userID idVO.UserID) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, exists := r.authentications[userID.String()]
	return exists, nil
}
