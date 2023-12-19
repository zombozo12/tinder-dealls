package resthttp

import (
	"context"
	"github.com/zombozo12/tinder-dealls/domain"
)

type AuthService interface {
	Login(ctx context.Context, req domain.AuthRequest) (*domain.AuthResponse, error)
	Register(ctx context.Context, req domain.AuthRequest) error
}

type ProfileService interface {
	Create(ctx context.Context, userID int64, req domain.ProfileRequest) error
	UpdateProfilePic(ctx context.Context, userID int64, req domain.UpdateProfilePicRequest) error
	UpdateProfile(ctx context.Context, userID int64, req domain.ProfileRequest) error
}

type MatcherService interface {
	Like(ctx context.Context, userID int64, targetUserID int64) error
	SuperLike(ctx context.Context, userID int64, targetUserID int64) error
	Dislike(ctx context.Context, userID int64, targetUserID int64) error
}

type RecommendationService interface {
	GetRecommendation(ctx context.Context, userID int64) (recommendation []domain.Profile, err error)
}
