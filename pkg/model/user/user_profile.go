package user

import (
	"context"
)

type UserProfile struct {
	PhoneNo string
}

func (up *UserProfile) Save(ctx context.Context) error {
	return nil
}
