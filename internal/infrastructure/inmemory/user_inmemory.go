package inmemory

import (
	"context"
	"sync"

	userDomain "github.com/u104rak1/pocgo/internal/domain/user"
	idVO "github.com/u104rak1/pocgo/internal/domain/value_object/id"
)

type userInMemoryRepository struct {
	mu    sync.RWMutex
	users map[string]*userDomain.User
}

func NewUserInMemoryRepository() userDomain.IUserRepository {
	return &userInMemoryRepository{
		users: make(map[string]*userDomain.User),
	}
}

func (r *userInMemoryRepository) Save(ctx context.Context, user *userDomain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.users[user.IDString()] = user
	return nil
}

func (r *userInMemoryRepository) FindByID(ctx context.Context, id idVO.UserID) (*userDomain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	user, exists := r.users[id.String()]
	if !exists {
		return nil, nil
	}
	return user, nil
}

func (r *userInMemoryRepository) FindByEmail(ctx context.Context, email string) (*userDomain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, user := range r.users {
		if user.Email() == email {
			return user, nil
		}
	}
	return nil, nil
}

func (r *userInMemoryRepository) ExistsByID(ctx context.Context, id idVO.UserID) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, exists := r.users[id.String()]
	return exists, nil
}

func (r *userInMemoryRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, user := range r.users {
		if user.Email() == email {
			return true, nil
		}
	}
	return false, nil
}
