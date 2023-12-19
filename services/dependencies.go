package services

import (
	"context"
	"github.com/zombozo12/tinder-dealls/domain"
)

type AuthRepoInterface interface {
	Login(ctx context.Context, req domain.AuthRequest) (user domain.User, err error)
	Register(ctx context.Context, req domain.AuthRequest) (domain.User, error)
	UpdateToken(ctx context.Context, userID int64, req domain.UpdateTokenRequest) error
}

type ProfileRepoInterface interface {
	Create(ctx context.Context, userId int64, req domain.ProfileRequest) error
	UpdateProfilePic(ctx context.Context, userID int64, req domain.UpdateProfilePicRequest) error
	UpdateProfile(ctx context.Context, userID int64, req domain.ProfileRequest) error
	GetProfile(ctx context.Context, userID int64) (profile *domain.Profile, err error)
	GetProfileRecommendation(ctx context.Context, interest string, notInUserID []int64) (profiles []domain.Profile, err error)
}

type RedisRepoInterface interface {
	Get(ctx context.Context, key string) (string, error)
	GetValues(ctx context.Context, key string) ([]string, error)
	Set(ctx context.Context, key string, value interface{}, expiration int) error
	Expire(ctx context.Context, key string, expiration int) error
	Incr(ctx context.Context, key string) error
	Exists(ctx context.Context, keys ...string) (bool, error)
}

type InventoryRepoInterface interface {
	Create(ctx context.Context, req domain.CreateInventoryRequest) error
	GetByUserId(ctx context.Context, userID int64) (inventory *domain.Inventory, err error)
	UpdateLikes(ctx context.Context, userID int64, likes int) error
	UpdateSuperLikes(ctx context.Context, userID int64, superLikes int) error
	UpdateSwipes(ctx context.Context, userID int64, swipes int) error
}

type MatchedRepoInterface interface {
	Create(ctx context.Context, req domain.MatchRequest) error
	IsMatched(ctx context.Context, req domain.MatchRequest) (bool, error)
	IsExists(ctx context.Context, req domain.MatchRequest) (bool, error)
}

type NotificationRepoInterface interface {
	Create(ctx context.Context, req domain.NotificationRequest) error
	GetAllByUserId(ctx context.Context, userID int64) ([]domain.Notification, error)
	SetRead(ctx context.Context, notificationID int64) error
}
