package matched

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
	create(ctx context.Context, req domain.MatchRequest) error
	isMatched(ctx context.Context, req domain.MatchRequest) (bool, error)
	isExists(ctx context.Context, req domain.MatchRequest) (bool, error)
}

func newDatabase(db *gorm.DB, cfg *domain.Config) dbInterface {
	return &dbModule{
		db:  db,
		cfg: cfg,
	}
}

func (d dbModule) create(ctx context.Context, req domain.MatchRequest) error {
	startTime := time.Now()
	tags := make(log.Fields)

	defer func() {
		tags["name"] = "repo.database.matched.create"
		tags["elapsed_time"] = time.Since(startTime).String()
		tags["request_id"] = ctx.Value("requestid").(string)
		log.WithFields(tags).Debug()
	}()

	match := domain.Matched{
		UserAID: req.UserID,
		UserBID: req.TargetUserID,
	}

	result := d.db.Table("matched").Create(&match)
	if result.Error != nil {
		tags["error"] = result.Error.Error()
		tags["status"] = "error"
		return result.Error
	}

	tags["status"] = "success"
	return nil
}

func (d dbModule) isMatched(ctx context.Context, req domain.MatchRequest) (bool, error) {
	startTime := time.Now()
	tags := make(log.Fields)

	defer func() {
		tags["name"] = "repo.database.matched.isMatched"
		tags["elapsed_time"] = time.Since(startTime).String()
		tags["request_id"] = ctx.Value("requestid").(string)
		log.WithFields(tags).Debug()
	}()

	var matched domain.Matched
	result := d.db.Table("matched").Where("user_a_id = ? AND user_b_id = ?", req.UserID, req.TargetUserID).First(&matched)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			tags["status"] = "not_found"
			return false, nil
		}

		tags["error"] = result.Error.Error()
		tags["status"] = "error"
		return false, result.Error
	}

	tags["status"] = "success"
	return true, nil
}

func (d dbModule) isExists(ctx context.Context, req domain.MatchRequest) (bool, error) {
	startTime := time.Now()
	tags := make(log.Fields)

	defer func() {
		tags["name"] = "repo.database.matched.isExists"
		tags["elapsed_time"] = time.Since(startTime).String()
		tags["request_id"] = ctx.Value("requestid").(string)
		log.WithFields(tags).Debug()
	}()

	var matched domain.Matched
	result := d.db.Table("matched").Where("user_a_id = ? AND user_b_id = ?", req.UserID, req.TargetUserID).First(&matched)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			tags["status"] = "not_found"
			return false, nil
		}

		tags["error"] = result.Error.Error()
		tags["status"] = "error"
		return false, result.Error
	}

	tags["status"] = "success"
	return true, nil
}
