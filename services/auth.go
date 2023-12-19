package services

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/goccy/go-json"
	"github.com/golang-jwt/jwt/v5"
	log "github.com/sirupsen/logrus"
	"github.com/zombozo12/tinder-dealls/domain"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"time"
)

type authServiceModule struct {
	cfg           *domain.Config
	db            *gorm.DB
	authRepo      AuthRepoInterface
	inventoryRepo InventoryRepoInterface
	profileRepo   ProfileRepoInterface
}

type AuthServiceInterface interface {
	Login(ctx context.Context, req domain.AuthRequest) (*domain.AuthResponse, error)
	Register(ctx context.Context, req domain.AuthRequest) error
}

func NewAuthService(cfg *domain.Config, db *gorm.DB, authRepo AuthRepoInterface, inventoryRepo InventoryRepoInterface,
	profileRepo ProfileRepoInterface) (AuthServiceInterface, error) {
	validate := validator.New()
	if err := validate.Struct(cfg); err != nil {
		return nil, err
	}

	return &authServiceModule{
		cfg:           cfg,
		db:            db,
		authRepo:      authRepo,
		inventoryRepo: inventoryRepo,
		profileRepo:   profileRepo,
	}, nil
}

//goland:noinspection GoLinter
func (a *authServiceModule) Login(ctx context.Context, req domain.AuthRequest) (*domain.AuthResponse, error) {
	startTime := time.Now()
	tags := make(log.Fields)
	defer func() {
		tags["name"] = "service.auth.login"
		tags["elapsed_time"] = time.Since(startTime).String()
		tags["request_id"] = ctx.Value("requestid").(string)
		log.WithFields(tags).Debug()
	}()

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		tags["error"] = "failed validating request"
		tags["status"] = "error"
		return nil, err
	}

	user, err := a.authRepo.Login(ctx, req)
	if err != nil {
		tags["error"] = "failed login"
		tags["status"] = "error"
		return nil, err
	}

	userJSON, err := json.Marshal(user)
	if err != nil {
		tags["error"] = "failed marshal user"
		tags["status"] = "error"
		return nil, err
	}

	userEncode, err := domain.EncryptAESWithGCM(string(userJSON), a.cfg.JWT.Key)
	if err != nil {
		tags["error"] = "failed encrypt user"
		tags["status"] = "error"
		tags["actual_error"] = err.Error()
		return nil, err
	}

	expireTime := time.Now().Add(time.Hour * time.Duration(a.cfg.JWT.ExpireTime)).Unix()
	claims := jwt.MapClaims{
		"sub": userEncode,
		"exp": expireTime,
		"iat": time.Now().Unix(),
		"nbf": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(a.cfg.Key))
	if err != nil {
		tags["error"] = "failed generating token"
		tags["status"] = "error"
		return nil, err
	}

	updateTokenReq := domain.UpdateTokenRequest{
		AccessToken:    tokenString,
		TokenExpiredAt: time.Unix(expireTime, 0),
	}

	err = a.authRepo.UpdateToken(ctx, user.ID, updateTokenReq)
	if err != nil {
		tags["error"] = "failed update token"
		tags["status"] = "error"
		return nil, err
	}

	response := &domain.AuthResponse{
		Token: tokenString,
	}

	tags["status"] = "success"
	return response, nil
}

func (a *authServiceModule) Register(ctx context.Context, req domain.AuthRequest) error {
	startTime := time.Now()
	tags := make(log.Fields)
	defer func() {
		tags["name"] = "service.auth.register"
		tags["elapsed_time"] = time.Since(startTime).String()
		tags["request_id"] = ctx.Value("requestid").(string)
		log.WithFields(tags).Debug()
	}()

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		tags["error"] = "failed validating request"
		tags["status"] = "error"
		return nil
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		tags["error"] = "failed hashing password"
		tags["status"] = "error"
		return err
	}

	req.Password = string(hashedPassword)

	user, err := a.authRepo.Register(ctx, req)
	if err != nil {
		tags["error"] = "failed register"
		tags["status"] = "error"
		return err
	}

	if err := a.inventoryRepo.Create(ctx, domain.CreateInventoryRequest{
		UserID:     user.ID,
		Likes:      10,
		SuperLikes: 1,
	}); err != nil {
		tags["error"] = "failed create inventory"
		tags["status"] = "error"
		return err
	}

	if err := a.profileRepo.Create(ctx, user.ID, domain.ProfileRequest{
		Name:       "",
		Gender:     "",
		InterestIn: "",
	}); err != nil {
		tags["error"] = "failed create profile"
		tags["status"] = "error"
		return err
	}

	tags["status"] = "success"
	return nil
}
