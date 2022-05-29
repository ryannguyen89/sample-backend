package memory

import (
	"context"
	"sync"

	"sampleBackend/internal/storage"
	"sampleBackend/internal/user"
)

type UserStorage struct {
	mu    sync.Mutex
	users map[string]user.User
}

func NewUserStorage() *UserStorage {
	return &UserStorage{
		users: make(map[string]user.User),
	}
}

func (us *UserStorage) Create(_ context.Context, u user.User) error {
	us.mu.Lock()
	defer us.mu.Unlock()

	if _, exist := us.users[u.Email]; exist {
		return storage.ErrAlreadyExist
	}

	us.users[u.Email] = u
	return nil
}

func (us *UserStorage) Verify(_ context.Context, u user.User) error {
	us.mu.Lock()
	defer us.mu.Unlock()

	if item, exist := us.users[u.Email]; exist {
		if item.Password == u.Password {
			return nil
		}
		return storage.ErrInvalidInfo
	} else {
		return storage.ErrInvalidInfo
	}
}
