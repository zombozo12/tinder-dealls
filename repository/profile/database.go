package profile

import (
	"context"
	"errors"
	log "github.com/sirupsen/logrus"
	"github.com/zombozo12/tinder-dealls/domain"
	"gorm.io/gorm"
	"time"
)

type module struct {
	db  *gorm.DB
	cfg *domain.Config
}

type dbInterface interface {
	create(ctx context.Context, userId int64, req domain.ProfileRequest) error
	updateProfilePic(ctx context.Context, userID int64, req domain.UpdateProfilePicRequest) error
	updateProfile(ctx context.Context, userID int64, req domain.ProfileRequest) error
	getProfile(ctx context.Context, userID int64) (profile *domain.Profile, err error)
	getProfileRecommendation(ctx context.Context, interest string, notInUserID []int64) (profiles []domain.Profile, err error)
}

func newDatabase(db *gorm.DB, cfg *domain.Config) dbInterface {
	return &module{
		db:  db,
		cfg: cfg,
	}
}

func (m module) create(ctx context.Context, userID int64, req domain.ProfileRequest) error {
	startTime := time.Now()
	tags := make(log.Fields)

	defer func() {
		tags["name"] = "repo.database.profile.create"
		tags["elapsed_time"] = time.Since(startTime).String()
		tags["request_id"] = ctx.Value("requestid").(string)
		log.WithFields(tags).Debug()
	}()

	profile := domain.Profile{
		UserID:     userID,
		Name:       req.Name,
		Gender:     req.Gender,
		InterestIn: req.InterestIn,
	}

	if result := m.db.Table("profile").Create(&profile); result.Error != nil {
		tags["error"] = result.Error.Error()
		tags["status"] = "error"
		return result.Error
	}

	tags["status"] = "success"
	return nil
}

func (m module) updateProfilePic(ctx context.Context, userID int64, req domain.UpdateProfilePicRequest) error {
	startTime := time.Now()
	tags := make(log.Fields)

	defer func() {
		tags["name"] = "repo.database.profile.update_profile_pic"
		tags["elapsed_time"] = time.Since(startTime).String()
		tags["request_id"] = ctx.Value("requestid").(string)
		log.WithFields(tags).Debug()
	}()

	if result := m.db.Table("profile").Where("user_id = ?", userID).Update("pic", req.Pic); result.Error != nil {
		tags["error"] = result.Error.Error()
		tags["status"] = "error"
		return result.Error
	}

	tags["status"] = "success"
	return nil
}

func (m module) updateProfile(ctx context.Context, userID int64, req domain.ProfileRequest) error {
	startTime := time.Now()
	tags := make(log.Fields)

	defer func() {
		tags["name"] = "repo.database.profile.update_profile"
		tags["elapsed_time"] = time.Since(startTime).String()
		tags["request_id"] = ctx.Value("requestid").(string)
		log.WithFields(tags).Debug()
	}()

	if result := m.db.Table("profile").Where("user_id = ?", userID).Update("name", req.Name); result.Error != nil {
		tags["error"] = result.Error.Error()
		tags["status"] = "error"
		return result.Error
	}

	tags["status"] = "success"
	return nil
}

func (m module) getProfile(ctx context.Context, userID int64) (profile *domain.Profile, err error) {
	startTime := time.Now()
	tags := make(log.Fields)

	defer func() {
		tags["name"] = "repo.database.profile.get_profile"
		tags["elapsed_time"] = time.Since(startTime).String()
		tags["request_id"] = ctx.Value("requestid").(string)
		tags["status"] = "success"
		log.WithFields(tags).Debug()
	}()

	if result := m.db.Table("profile").Where("user_id = ?", userID).First(&profile); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			tags["status"] = "not_found"
			return nil, nil
		}

		tags["error"] = result.Error.Error()
		tags["status"] = "error"
		return nil, result.Error
	}

	tags["status"] = "success"
	return profile, nil
}

func (m module) getProfileRecommendation(ctx context.Context, interest string, notInUserID []int64) (profiles []domain.Profile, err error) {
	startTime := time.Now()
	tags := make(log.Fields)

	defer func() {
		tags["name"] = "repo.database.profile.get_profile_recommendation"
		tags["elapsed_time"] = time.Since(startTime).String()
		tags["request_id"] = ctx.Value("requestid").(string)
		tags["status"] = "success"
		log.WithFields(tags).Debug()
	}()

	if result := m.db.Table("profile").
		Where("interest_in = ? AND gender = ? AND user_id NOT IN (?)", interest, interest, notInUserID).
		Limit(10).
		Find(&profiles); result.Error != nil {
		tags["error"] = result.Error.Error()
		tags["status"] = "error"
		return profiles, result.Error
	}

	tags["status"] = "success"
	return profiles, nil
}
