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
	return []*user.User{
		{ID: 1},
		{ID: 2},
		{ID: 3},
		{ID: 4},
		{ID: 5},
		{ID: 6},
		{ID: 7},
		{ID: 8},
		{ID: 9},
		{ID: 10},
		{ID: 11},
		{ID: 12},
		{ID: 13},
		{ID: 14},
		{ID: 15},
		{ID: 16},
		{ID: 17},
		{ID: 18},
		{ID: 19},
		{ID: 20},
	}, nil
}
