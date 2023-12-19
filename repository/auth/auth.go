package auth

import (
	"context"
	"github.com/zombozo12/tinder-dealls/domain"
	"gorm.io/gorm"
)

type Module struct {
	dbs dbInterface
	cfg *domain.Config
}

func New(db *gorm.DB, cfg *domain.Config) *Module {
	return &Module{
		dbs: newDatabase(db, cfg),
		cfg: cfg,
	}
}

func (m Module) Login(ctx context.Context, req domain.AuthRequest) (user domain.User, err error) {
	return m.dbs.login(ctx, req)
}

func (m Module) UpdateToken(ctx context.Context, userID int64, req domain.UpdateTokenRequest) error {
	return m.dbs.updateToken(ctx, userID, req)
}

func (m Module) Register(ctx context.Context, req domain.AuthRequest) (domain.User, error) {
	return m.dbs.register(ctx, req)
}
