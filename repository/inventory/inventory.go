package inventory

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

func (m Module) Create(ctx context.Context, req domain.CreateInventoryRequest) error {
	return m.dbs.create(ctx, req)
}

func (m Module) GetByUserId(ctx context.Context, userID int64) (inventory *domain.Inventory, err error) {
	return m.dbs.getByUserId(ctx, userID)
}

func (m Module) UpdateLikes(ctx context.Context, userID int64, likes int) error {
	return m.dbs.updateLikes(ctx, userID, likes)
}

func (m Module) UpdateSuperLikes(ctx context.Context, userID int64, superLikes int) error {
	return m.dbs.updateSuperLikes(ctx, userID, superLikes)
}

func (m Module) UpdateSwipes(ctx context.Context, userID int64, swipes int) error {
	return m.dbs.updateSwipes(ctx, userID, swipes)
}
