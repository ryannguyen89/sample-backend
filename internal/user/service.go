package user

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"

	"sampleBackend/internal/storage"
)

var (
	ErrUserExist   = errors.New("user exist")
	ErrUserInvalid = errors.New("user invalid")

	secretKey = []byte("G+KbPeSh")
)

type Storage interface {
	Create(ctx context.Context, u User) error
	Verify(ctx context.Context, u User) error
}

type Service struct {
	storage Storage
}

func NewService(s Storage) *Service {
	return &Service{storage: s}
}

func (s *Service) CreateUser(ctx context.Context, u User) error {
	err := s.storage.Create(ctx, u)
	if err != nil {
		if storage.IsErrAlreadyExist(err) {
			return fmt.Errorf("create user: %v - %w", err, ErrUserExist)
		}
		return fmt.Errorf("create user: %w", err)
	}
	return nil
}

func (s *Service) Login(ctx context.Context, u User) (*Login, error) {
	// Check user
	err := s.storage.Verify(ctx, u)
	if err != nil {
		if storage.IsErrInvalidInfo(err) {
			return nil, fmt.Errorf("verify user: %v - %w", err, ErrUserInvalid)
		}
		return nil, fmt.Errorf("verify user: %w", err)
	}

	// Generate token
	t := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Audience: jwt.ClaimStrings{u.Email},
		IssuedAt: &jwt.NumericDate{
			Time: t,
		},
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return nil, fmt.Errorf("sign string: %w", err)
	}

	return &Login{
		Token: tokenString,
	}, nil
}

func (s *Service) ValidateToken(_ context.Context, tokenString string) error {
	_, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return fmt.Errorf("parse token: %w", err)
	}

	return nil
}

func IsErrUserExist(err error) bool {
	return errors.Is(err, ErrUserExist)
}

func IsErrUserInvalid(err error) bool {
	return errors.Is(err, ErrUserInvalid)
}
