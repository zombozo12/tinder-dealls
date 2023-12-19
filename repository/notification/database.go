package notification

import (
	"context"
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
	create(ctx context.Context, req domain.NotificationRequest) error
	getAllByUserId(ctx context.Context, userID int64) ([]domain.Notification, error)
	setRead(ctx context.Context, notificationID int64) error
}

func newDatabase(db *gorm.DB, cfg *domain.Config) dbInterface {
	return &dbModule{
		db:  db,
		cfg: cfg,
	}
}

func (d dbModule) create(ctx context.Context, req domain.NotificationRequest) error {
	startTime := time.Now()
	tags := make(log.Fields)

	defer func() {
		tags["name"] = "repo.database.notification.create"
		tags["elapsed_time"] = time.Since(startTime).String()
		tags["request_id"] = ctx.Value("requestid").(string)
		log.WithFields(tags).Debug()
	}()

	notification := domain.Notification{
		UserID:  req.UserID,
		Message: req.Message,
		IsRead:  false,
	}

	result := d.db.Table("notification").Create(&notification)
	if result.Error != nil {
		tags["error"] = result.Error.Error()
		tags["status"] = "error"
		return result.Error
	}

	tags["status"] = "success"
	return nil
}

func (d dbModule) getAllByUserId(ctx context.Context, userID int64) ([]domain.Notification, error) {
	startTime := time.Now()
	tags := make(log.Fields)

	defer func() {
		tags["name"] = "repo.database.notification.getAll"
		tags["elapsed_time"] = time.Since(startTime).String()
		tags["request_id"] = ctx.Value("requestid").(string)
		log.WithFields(tags).Debug()
	}()

	var notifications []domain.Notification
	result := d.db.Table("notification").
		Where("user_id = ?", userID).
		Find(&notifications)
	if result.Error != nil {
		tags["error"] = result.Error.Error()
		tags["status"] = "error"
		return nil, result.Error
	}

	tags["status"] = "success"
	return notifications, nil
}

func (d dbModule) setRead(ctx context.Context, notificationID int64) error {
	startTime := time.Now()
	tags := make(log.Fields)

	defer func() {
		tags["name"] = "repo.database.notification.setRead"
		tags["elapsed_time"] = time.Since(startTime).String()
		tags["request_id"] = ctx.Value("requestid").(string)
		log.WithFields(tags).Debug()
	}()

	result := d.db.Table("notification").
		Where("id = ?", notificationID).
		Update("is_read", true)
	if result.Error != nil {
		tags["error"] = result.Error.Error()
		tags["status"] = "error"
		return result.Error
	}

	tags["status"] = "success"
	return nil
}
