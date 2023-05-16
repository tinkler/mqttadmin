package page

import (
	"context"

	"github.com/tinkler/mqttadmin/pkg/model/user"
)

type Page struct {
	Page    int
	PerPage int
	Total   int
}

func (p *Page) FetchUser(ctx context.Context) ([]*user.User, error) {
	return []*user.User{}, nil
}
