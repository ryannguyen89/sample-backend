package memory

import (
	"context"
	"sync"

	"sampleBackend/internal/storage"
	"sampleBackend/internal/user"
)

type UserStorage struct {
	mu    sync.Mutex
	users []user.User
}

func NewUserStorage() *UserStorage {
	return &UserStorage{}
}

func (us *UserStorage) Create(_ context.Context, u user.User) error {
	us.mu.Lock()
	defer us.mu.Unlock()

	for _, i := range us.users {
		if i.Email == u.Email {
			return storage.ErrAlreadyExist
		}
	}

	us.users = append(us.users, u)
	return nil
}

func (us *UserStorage) Verify(ctx context.Context, u user.User) (err error) {
	us.mu.Lock()
	defer us.mu.Unlock()

	err = storage.ErrInvalidInfo

	for _, i := range us.users {
		if i.Email == u.Email {
			if i.Password == u.Password {
				err = nil
			}
			break
		}
	}

	return
}
