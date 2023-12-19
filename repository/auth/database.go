package auth

import (
	"context"
	"errors"
	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
	"github.com/zombozo12/tinder-dealls/domain"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"time"
)

type module struct {
	db  *gorm.DB
	cfg *domain.Config
}

type dbInterface interface {
	login(ctx context.Context, req domain.AuthRequest) (user domain.User, err error)
	register(ctx context.Context, req domain.AuthRequest) (domain.User, error)
	updateToken(ctx context.Context, userID int64, req domain.UpdateTokenRequest) error
}

func newDatabase(db *gorm.DB, cfg *domain.Config) dbInterface {
	return &module{
		db:  db,
		cfg: cfg,
	}
}

func (m module) login(ctx context.Context, req domain.AuthRequest) (user domain.User, err error) {
	startTime := time.Now()
	tags := make(log.Fields)

	defer func() {
		tags["name"] = "repo.database.auth.login"
		tags["elapsed_time"] = time.Since(startTime).String()
		tags["request_id"] = ctx.Value("requestid").(string)
		log.WithFields(tags).Debug()
	}()

	if result := m.db.Table("users").Where("email = ?", req.Email).First(&user); result.Error != nil {
		tags["error"] = result.Error.Error()
		tags["status"] = "error"
		return user, result.Error
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		tags["error"] = err.Error()
		tags["status"] = "error"
		return user, err
	}

	user = domain.User{
		ID:        user.ID,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	tags["status"] = "success"
	return user, nil
}

func (m module) register(ctx context.Context, req domain.AuthRequest) (domain.User, error) {
	startTime := time.Now()
	tags := make(log.Fields)

	defer func() {
		tags["name"] = "repo.database.auth.register"
		tags["elapsed_time"] = time.Since(startTime).String()
		tags["request_id"] = ctx.Value("requestid").(string)
		log.WithFields(tags).Debug()
	}()

	user := domain.User{
		Email:    req.Email,
		Password: req.Password,
	}
	result := m.db.Table("users").Create(&user)
	if result.Error != nil {
		tags["error"] = "failed creating user"
		tags["status"] = "error"
		return user, result.Error
	}

	tags["status"] = "success"
	return user, nil
}

func (m module) updateToken(ctx context.Context, userID int64, req domain.UpdateTokenRequest) error {
	startTime := time.Now()
	tags := make(log.Fields)

	defer func() {
		tags["name"] = "repo.database.auth.update_token"
		tags["elapsed_time"] = time.Since(startTime).String()
		tags["request_id"] = ctx.Value("requestid").(string)
		log.WithFields(tags).Debug()
	}()

	if userID == 0 {
		tags["error"] = "id is empty"
		tags["status"] = "error"
		return errors.New("id is empty")
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		tags["error"] = "failed validating request"
		tags["status"] = "error"
		return err
	}

	req.UpdatedAt = time.Now()

	result := m.db.Table("users").Where("id = ?", userID).
		Updates(req)
	if result.Error != nil {
		tags["error"] = result.Error.Error()
		tags["status"] = "error"
		return result.Error
	}

	tags["status"] = "success"
	return nil
}
