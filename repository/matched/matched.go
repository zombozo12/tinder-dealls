package matched

import (
	"context"
	"github.com/zombozo12/tinder-dealls/domain"
	"gorm.io/gorm"
)

type Module struct {
	cfg *domain.Config
	dbs dbInterface
}

func New(db *gorm.DB, cfg *domain.Config) *Module {
	return &Module{
		cfg: cfg,
		dbs: newDatabase(db, cfg),
	}
}

func (m Module) Create(ctx context.Context, req domain.MatchRequest) error {
	return m.dbs.create(ctx, req)
}

func (m Module) IsMatched(ctx context.Context, req domain.MatchRequest) (bool, error) {
	return m.dbs.isMatched(ctx, req)
}

func (m Module) IsExists(ctx context.Context, req domain.MatchRequest) (bool, error) {
	return m.dbs.isExists(ctx, req)
}
