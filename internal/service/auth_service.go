package service

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mysecodgit/go_accounting/internal/dto"
	"github.com/mysecodgit/go_accounting/internal/store"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

type AuthUserStore interface {
	GetByUsername(ctx context.Context, username string) (*store.User, error)
}

type AuthService struct {
	userStore AuthUserStore
	jwtSecret []byte
	tokenTTL  time.Duration
}

func NewAuthService(userStore AuthUserStore, jwtSecret string) *AuthService {
	return &AuthService{
		userStore: userStore,
		jwtSecret: []byte(jwtSecret),
		tokenTTL:  24 * time.Hour,
	}
}

func (s *AuthService) Login(ctx context.Context, username string, password string) (*dto.LoginResponse, error) {
	u, err := s.userStore.GetByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	// NOTE: Your current DB schema stores plaintext passwords (varchar(10)).
	// This is not recommended, but we validate it as-is to match your schema.
	// TODO: Hash the password before storing it in the database
	if u.Password != password {
		return nil, ErrInvalidCredentials
	}

	now := time.Now()
	claims := jwt.MapClaims{
		"sub":      u.ID,
		"username": u.Username,
		"iat":      now.Unix(),
		"exp":      now.Add(s.tokenTTL).Unix(),
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(s.jwtSecret)
	if err != nil {
		return nil, err
	}

	// Never return password
	u.Password = ""

	return &dto.LoginResponse{
		AccessToken: token,
		Username:    u.Username,
		User:        *u,
	}, nil
}
