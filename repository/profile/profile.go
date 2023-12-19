package profile

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

func (m Module) Create(ctx context.Context, userId int64, req domain.ProfileRequest) error {
	return m.dbs.create(ctx, userId, req)
}

func (m Module) UpdateProfilePic(ctx context.Context, userID int64, req domain.UpdateProfilePicRequest) error {
	return m.dbs.updateProfilePic(ctx, userID, req)
}

func (m Module) UpdateProfile(ctx context.Context, userID int64, req domain.ProfileRequest) error {
	return m.dbs.updateProfile(ctx, userID, req)
}

func (m Module) GetProfile(ctx context.Context, userID int64) (profile *domain.Profile, err error) {
	return m.dbs.getProfile(ctx, userID)
}

func (m Module) GetProfileRecommendation(ctx context.Context, interest string, notInUserID []int64) (profiles []domain.Profile, err error) {
	return m.dbs.getProfileRecommendation(ctx, interest, notInUserID)
}
