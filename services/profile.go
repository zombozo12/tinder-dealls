package services

import (
	"context"
	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
	"github.com/zombozo12/tinder-dealls/domain"
	"gorm.io/gorm"
	"time"
)

type profileServiceModule struct {
	cfg         *domain.Config
	db          *gorm.DB
	profileRepo ProfileRepoInterface
}

type ProfileServiceModuleInterface interface {
	Create(ctx context.Context, userId int64, req domain.ProfileRequest) error
	UpdateProfilePic(ctx context.Context, userID int64, req domain.UpdateProfilePicRequest) error
	UpdateProfile(ctx context.Context, userID int64, req domain.ProfileRequest) error
}

func NewProfileService(cfg *domain.Config, db *gorm.DB, profileRepo ProfileRepoInterface) (ProfileServiceModuleInterface, error) {
	return &profileServiceModule{
		cfg:         cfg,
		db:          db,
		profileRepo: profileRepo,
	}, nil
}

func (p profileServiceModule) Create(ctx context.Context, userId int64, req domain.ProfileRequest) error {
	startTime := time.Now()
	tags := make(log.Fields)

	defer func() {
		tags["name"] = "service.profile.create"
		tags["elapsed_time"] = time.Since(startTime).String()
		tags["request_id"] = ctx.Value("requestid").(string)
		log.WithFields(tags).Debug()
	}()

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		tags["error"] = "failed validating request"
		tags["status"] = "error"
		return err
	}

	if err := p.profileRepo.Create(ctx, userId, req); err != nil {
		tags["error"] = "failed to create profile"
		tags["status"] = "error"
		return err
	}

	tags["status"] = "success"
	return nil
}

func (p profileServiceModule) UpdateProfilePic(ctx context.Context, userID int64, req domain.UpdateProfilePicRequest) error {
	startTime := time.Now()
	tags := make(log.Fields)

	defer func() {
		tags["name"] = "service.profile.update_profile_pic"
		tags["elapsed_time"] = time.Since(startTime).String()
		tags["request_id"] = ctx.Value("requestid").(string)
		log.WithFields(tags).Debug()
	}()

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		tags["error"] = "failed validating request"
		tags["status"] = "error"
		return err
	}

	if err := p.profileRepo.UpdateProfilePic(ctx, userID, req); err != nil {
		tags["error"] = "failed to update profile pic"
		tags["status"] = "error"
		return err
	}

	tags["status"] = "success"
	return nil
}

func (p profileServiceModule) UpdateProfile(ctx context.Context, userID int64, req domain.ProfileRequest) error {
	startTime := time.Now()
	tags := make(log.Fields)

	defer func() {
		tags["name"] = "service.profile.update_profile"
		tags["elapsed_time"] = time.Since(startTime).String()
		tags["request_id"] = ctx.Value("requestid").(string)
		log.WithFields(tags).Debug()
	}()

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		tags["error"] = "failed validating request"
		tags["status"] = "error"
		return err
	}

	if err := p.profileRepo.UpdateProfile(ctx, userID, req); err != nil {
		tags["error"] = "failed to update profile"
		tags["status"] = "error"
		return err
	}

	tags["status"] = "success"
	return nil
}
