package notification

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

func (m Module) Create(ctx context.Context, req domain.NotificationRequest) error {
	return m.dbs.create(ctx, req)
}

func (m Module) GetAllByUserId(ctx context.Context, userID int64) ([]domain.Notification, error) {
	return m.dbs.getAllByUserId(ctx, userID)
}

func (m Module) SetRead(ctx context.Context, notificationID int64) error {
	return m.dbs.setRead(ctx, notificationID)
}
