package services

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"
	"github.com/zombozo12/tinder-dealls/domain"
	"gorm.io/gorm"
	"time"
)

type recommendationServiceModule struct {
	cfg         *domain.Config
	db          *gorm.DB
	profileRepo ProfileRepoInterface
	redisRepo   RedisRepoInterface
}

type RecommendationServiceInterface interface {
	GetRecommendation(ctx context.Context, userID int64) (recommendation []domain.Profile, err error)
}

func NewRecommendationService(cfg *domain.Config, db *gorm.DB, redisRepo RedisRepoInterface,
	profileRepo ProfileRepoInterface) (RecommendationServiceInterface, error) {
	return &recommendationServiceModule{
		cfg:         cfg,
		db:          db,
		profileRepo: profileRepo,
		redisRepo:   redisRepo,
	}, nil
}

func (r recommendationServiceModule) GetRecommendation(ctx context.Context, userID int64) (recommendation []domain.Profile, err error) {
	startTime := time.Now()
	tags := make(log.Fields)

	defer func() {
		tags["name"] = "service.recommendation.get_recommendation"
		tags["elapsed_time"] = time.Since(startTime).String()
		tags["request_id"] = ctx.Value("requestid").(string)
		log.WithFields(tags).Debug()
	}()

	profile, err := r.profileRepo.GetProfile(ctx, userID)
	if err != nil {
		tags["error"] = "failed get profile"
		tags["status"] = "error"
		return nil, err
	}

	key := fmt.Sprintf("frozen:%d", profile.UserID)

	frozenIDs, err := r.redisRepo.Get(ctx, key)
	if err != nil {
		tags["error"] = "failed get frozen ids"
		tags["status"] = "error"
		return nil, err
	}

	var unmarshalFrozenIds []int64
	if frozenIDs != "" {
		if unmarshalErr := json.Unmarshal([]byte(frozenIDs), &unmarshalFrozenIds); unmarshalErr != nil {
			tags["error"] = "failed unmarshal"
			tags["status"] = "error"
			return nil, unmarshalErr
		}
	}

	var excludedIDs []int64
	if len(unmarshalFrozenIds) > 0 {
		excludedIDs = append(excludedIDs, unmarshalFrozenIds...)
	}

	if !lo.Contains(excludedIDs, userID) {
		excludedIDs = append(excludedIDs, userID)
	}

	profileRecommendations, err := r.profileRepo.GetProfileRecommendation(ctx, profile.InterestIn, excludedIDs)
	if err != nil {
		tags["error"] = "failed get recommendation"
		tags["status"] = "error"
		return nil, err
	}

	for _, v := range profileRecommendations {
		excludedIDs = append(excludedIDs, v.UserID)
	}

	marshaledFrozenIds, err := json.Marshal(excludedIDs)
	if err != nil {
		tags["error"] = "failed marshal"
		tags["status"] = "error"
		return nil, err
	}

	if err := r.redisRepo.Set(ctx, key, marshaledFrozenIds, 60*60*24); err != nil {
		tags["error"] = "failed set frozen ids"
		tags["status"] = "error"
		return nil, err
	}

	tags["status"] = "success"
	return profileRecommendations, nil
}
