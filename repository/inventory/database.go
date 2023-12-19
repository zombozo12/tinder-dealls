package inventory

import (
	"context"
	"errors"
	log "github.com/sirupsen/logrus"
	"github.com/zombozo12/tinder-dealls/domain"
	"gorm.io/gorm"
	"time"
)

type dbModule struct {
	db  *gorm.DB
	cfg *domain.Config
}

type dbInterface interface {
	create(ctx context.Context, req domain.CreateInventoryRequest) error
	getByUserId(ctx context.Context, userID int64) (inventory *domain.Inventory, err error)
	updateLikes(ctx context.Context, userID int64, likes int) error
	updateSuperLikes(ctx context.Context, userID int64, superLikes int) error
	updateSwipes(ctx context.Context, userID int64, swipes int) error
}

func newDatabase(db *gorm.DB, cfg *domain.Config) dbInterface {
	return &dbModule{
		db:  db,
		cfg: cfg,
	}
}

func (d dbModule) create(ctx context.Context, req domain.CreateInventoryRequest) error {
	startTime := time.Now()
	tags := make(log.Fields)

	defer func() {
		tags["name"] = "repo.database.inventory.create"
		tags["elapsed_time"] = time.Since(startTime).String()
		tags["request_id"] = ctx.Value("requestid").(string)
		log.WithFields(tags).Debug()
	}()

	inventory := domain.Inventory{
		UserID:     req.UserID,
		Likes:      req.Likes,
		SuperLikes: req.SuperLikes,
	}

	result := d.db.Table("inventory").Create(&inventory)
	if result.Error != nil {
		tags["error"] = result.Error.Error()
		tags["status"] = "error"
		return result.Error
	}

	tags["status"] = "success"
	return nil
}

func (d dbModule) getByUserId(ctx context.Context, userID int64) (inventory *domain.Inventory, err error) {
	startTime := time.Now()
	tags := make(log.Fields)

	defer func() {
		tags["name"] = "repo.database.inventory.get"
		tags["elapsed_time"] = time.Since(startTime).String()
		tags["request_id"] = ctx.Value("requestid").(string)
		log.WithFields(tags).Debug()
	}()

	result := d.db.Table("inventory").Where("user_id = ?", userID).First(&inventory)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			tags["status"] = "not_found"
			return nil, nil
		}

		tags["error"] = result.Error.Error()
		tags["status"] = "error"
		return nil, result.Error
	}

	tags["status"] = "success"
	return inventory, nil
}

func (d dbModule) updateLikes(ctx context.Context, userID int64, likes int) error {
	startTime := time.Now()
	tags := make(log.Fields)

	defer func() {
		tags["name"] = "repo.database.inventory.update_likes"
		tags["elapsed_time"] = time.Since(startTime).String()
		tags["request_id"] = ctx.Value("requestid").(string)
		log.WithFields(tags).Debug()
	}()

	result := d.db.Table("inventory").Where("user_id = ?", userID).Update("likes", likes)
	if result.Error != nil {
		tags["error"] = result.Error.Error()
		tags["status"] = "error"
		return result.Error
	}

	tags["status"] = "success"
	return nil
}

func (d dbModule) updateSuperLikes(ctx context.Context, userID int64, superLikes int) error {
	startTime := time.Now()
	tags := make(log.Fields)

	defer func() {
		tags["name"] = "repo.database.inventory.update_super_likes"
		tags["elapsed_time"] = time.Since(startTime).String()
		tags["request_id"] = ctx.Value("requestid").(string)
		log.WithFields(tags).Debug()
	}()

	result := d.db.Table("inventory").Where("user_id = ?", userID).Update("super_likes", superLikes)
	if result.Error != nil {
		tags["error"] = result.Error.Error()
		tags["status"] = "error"
		return result.Error
	}

	tags["status"] = "success"
	return nil
}

func (d dbModule) updateSwipes(ctx context.Context, userID int64, swipes int) error {
	startTime := time.Now()
	tags := make(log.Fields)

	defer func() {
		tags["name"] = "repo.database.inventory.update_swipes"
		tags["elapsed_time"] = time.Since(startTime).String()
		tags["request_id"] = ctx.Value("requestid").(string)
		log.WithFields(tags).Debug()
	}()

	result := d.db.Table("inventory").Where("user_id = ?", userID).Update("swipes", swipes)
	if result.Error != nil {
		tags["error"] = result.Error.Error()
		tags["status"] = "error"
		return result.Error
	}

	tags["status"] = "success"
	return nil
}
