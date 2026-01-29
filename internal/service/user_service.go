package service

import (
	"context"

	"github.com/mysecodgit/go_accounting/internal/store"
)

// Interface lives here
type UserStore interface {
	GetAll(ctx context.Context) ([]store.User, error)
	GetAllByParentID(ctx context.Context, parentUserID int64) ([]store.User, error)
	GetByID(ctx context.Context, id int64) (*store.User, error)
	Create(ctx context.Context, user *store.User) error
	Update(ctx context.Context, user *store.User) error
	Delete(ctx context.Context, id int64) error
}

type UserService struct {
	userStore UserStore
}

func NewUserService(userStore UserStore) *UserService {
	return &UserService{userStore: userStore}
}

func (s *UserService) GetAll(ctx context.Context) ([]store.User, error) {
	return s.userStore.GetAll(ctx)
}

func (s *UserService) GetAllByParentID(ctx context.Context, parentUserID int64) ([]store.User, error) {
	return s.userStore.GetAllByParentID(ctx, parentUserID)
}

func (s *UserService) GetByID(ctx context.Context, id int64) (*store.User, error) {
	return s.userStore.GetByID(ctx, id)
}

func (s *UserService) Create(ctx context.Context, user *store.User) error {
	return s.userStore.Create(ctx, user)
}

func (s *UserService) Update(ctx context.Context, user *store.User) error {
	return s.userStore.Update(ctx, user)
}

func (s *UserService) Delete(ctx context.Context, id int64) error {
	return s.userStore.Delete(ctx, id)
}
