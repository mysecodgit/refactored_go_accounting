package dto

import "github.com/mysecodgit/go_accounting/internal/store"

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	AccessToken string     `json:"accessToken"`
	Username    string     `json:"username"`
	User        store.User `json:"user"`
}
