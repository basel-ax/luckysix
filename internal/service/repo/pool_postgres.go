package repo

import (
	"github.com/basel-ax/luckysix/pkg/postgres"
)

// PoolRepo -.
type PoolRepo struct {
	*postgres.Postgres
}

// PoolRepo -.
func NewPoolRepo(pg *postgres.Postgres) *PoolRepo {
	return &PoolRepo{pg}
}
