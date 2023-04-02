package user

import (
	"context"
)

type Auth struct {
	Username string
	Password string
}

func (a *Auth) Login(ctx context.Context) (*User, error) {

	return &User{ID: 12}, nil
}
