package services

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	log "github.com/sirupsen/logrus"
	"github.com/zombozo12/tinder-dealls/domain"
	"gorm.io/gorm"
	"reflect"
	"testing"
	"time"
)

func TestNewAuthService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		cfg           *domain.Config
		db            *gorm.DB
		authRepo      AuthRepoInterface
		inventoryRepo InventoryRepoInterface
		profileRepo   ProfileRepoInterface
	}

	config := &domain.Config{
		Server: domain.Server{
			Port: 3000,
		},
		Key: "asd",
		Database: domain.Database{
			Host:     "test",
			Port:     1234,
			Username: "test",
			Password: "test",
			Name:     "test",
		},
		JWT: domain.JWT{
			Key:        "asd",
			ExpireTime: 30,
		},
		Redis: domain.Redis{
			Host:     "test",
			Port:     1234,
			Password: "",
		},
	}

	failedConfig := &domain.Config{}

	db := &gorm.DB{}

	authMock := NewMockAuthRepoInterface(ctrl)
	inventoryMock := NewMockInventoryRepoInterface(ctrl)
	profileMock := NewMockProfileRepoInterface(ctrl)

	tests := []struct {
		name    string
		args    args
		want    AuthServiceInterface
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				cfg:           config,
				db:            db,
				authRepo:      authMock,
				inventoryRepo: inventoryMock,
				profileRepo:   profileMock,
			},
			want: &authServiceModule{
				cfg:           config,
				db:            db,
				authRepo:      authMock,
				inventoryRepo: inventoryMock,
				profileRepo:   profileMock,
			},
			wantErr: false,
		},
		{
			name: "failed",
			args: args{
				cfg:           failedConfig,
				db:            db,
				authRepo:      authMock,
				inventoryRepo: inventoryMock,
				profileRepo:   profileMock,
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewAuthService(tt.args.cfg, tt.args.db, tt.args.authRepo, tt.args.inventoryRepo, tt.args.profileRepo)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewAuthService() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAuthService() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_authServiceModule_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		ctx context.Context
		req domain.AuthRequest
	}

	ctx := context.WithValue(context.Background(), "requestid", "test")

	config := &domain.Config{
		Server: domain.Server{
			Port: 3000,
		},
		Key: "12345",
		Database: domain.Database{
			Host:     "test",
			Port:     1234,
			Username: "test",
			Password: "test",
			Name:     "test",
		},
		JWT: domain.JWT{
			Key:        "TRTgoAnzX&NPBAhA53C6PaMB&E5*d7wx",
			ExpireTime: 30,
		},
		Redis: domain.Redis{
			Host:     "test",
			Port:     1234,
			Password: "",
		},
	}

	failedConfig := &domain.Config{
		Server: domain.Server{
			Port: 3000,
		},
		Key: "12345",
		Database: domain.Database{
			Host:     "test",
			Port:     1234,
			Username: "test",
			Password: "test",
			Name:     "test",
		},
		JWT: domain.JWT{
			Key:        "test",
			ExpireTime: 30,
		},
		Redis: domain.Redis{
			Host:     "test",
			Port:     1234,
			Password: "",
		},
	}

	db := &gorm.DB{}

	tests := []struct {
		name    string
		args    args
		mock    func() AuthServiceInterface
		want    *domain.AuthResponse
		wantErr bool
	}{
		{
			name: "failed validating request",
			args: args{
				ctx: ctx,
				req: domain.AuthRequest{
					Email:    "",
					Password: "",
				},
			},
			mock: func() AuthServiceInterface {
				return &authServiceModule{
					cfg: failedConfig,
					db:  db,
				}
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "failed login",
			args: args{
				ctx: ctx,
				req: domain.AuthRequest{
					Email:    "test@mail.com",
					Password: "testtest",
				},
			},
			mock: func() AuthServiceInterface {
				authMock := NewMockAuthRepoInterface(ctrl)
				authMock.EXPECT().Login(ctx, domain.AuthRequest{
					Email:    "test@mail.com",
					Password: "testtest",
				}).Return(domain.User{}, gorm.ErrRecordNotFound)
				return &authServiceModule{
					cfg:      failedConfig,
					db:       db,
					authRepo: authMock,
				}
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "failed encrypt aes with gcm",
			args: args{
				ctx: ctx,
				req: domain.AuthRequest{
					Email:    "test@mail.com",
					Password: "testtest",
				},
			},
			mock: func() AuthServiceInterface {
				log.SetLevel(log.DebugLevel)
				authMock := NewMockAuthRepoInterface(ctrl)
				authMock.EXPECT().Login(ctx, domain.AuthRequest{
					Email:    "test@mail.com",
					Password: "testtest",
				}).Return(domain.User{
					ID:             1,
					Email:          "test@mail.com",
					Password:       "testtest",
					AccessToken:    "",
					TokenExpiredAt: nil,
					CreatedAt:      time.Time{},
					UpdatedAt:      time.Time{},
					DeletedAt:      nil,
				}, nil)

				return &authServiceModule{
					cfg:      failedConfig,
					db:       db,
					authRepo: authMock,
				}
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "failed update token",
			args: args{
				ctx: ctx,
				req: domain.AuthRequest{
					Email:    "test@mail.com",
					Password: "testtest",
				},
			},
			mock: func() AuthServiceInterface {
				authMock := NewMockAuthRepoInterface(ctrl)
				authMock.EXPECT().Login(ctx, domain.AuthRequest{
					Email:    "test@mail.com",
					Password: "testtest",
				}).Return(domain.User{
					ID:             1,
					Email:          "test@mail.com",
					Password:       "testtest",
					AccessToken:    "",
					TokenExpiredAt: nil,
					CreatedAt:      time.Time{},
					UpdatedAt:      time.Time{},
					DeletedAt:      nil,
				}, nil)

				authMock.EXPECT().UpdateToken(ctx, int64(1), gomock.Any()).Return(errors.New("test"))

				return &authServiceModule{
					cfg:      config,
					db:       db,
					authRepo: authMock,
				}
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "success",
			args: args{
				ctx: ctx,
				req: domain.AuthRequest{
					Email:    "test@mail.com",
					Password: "testtest",
				},
			},
			mock: func() AuthServiceInterface {
				authMock := NewMockAuthRepoInterface(ctrl)
				authMock.EXPECT().Login(ctx, domain.AuthRequest{
					Email:    "test@mail.com",
					Password: "testtest",
				}).Return(domain.User{
					ID:             1,
					Email:          "test@mail.com",
					Password:       "testtest",
					AccessToken:    "",
					TokenExpiredAt: nil,
					CreatedAt:      time.Time{},
					UpdatedAt:      time.Time{},
					DeletedAt:      nil,
				}, nil)

				authMock.EXPECT().UpdateToken(ctx, gomock.Any(), gomock.Any()).Return(nil)

				return &authServiceModule{
					cfg:      config,
					db:       db,
					authRepo: authMock,
				}
			},
			want: &domain.AuthResponse{
				Token: "", // gomock.Any() here
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := tt.mock()
			got, err := a.Login(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(tt.want, got) {
				tt.want.Token = got.Token
				if !reflect.DeepEqual(tt.want, got) {
					t.Errorf("Login() got = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func Test_authServiceModule_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		ctx context.Context
		req domain.AuthRequest
	}

	ctx := context.WithValue(context.Background(), "requestid", "test")

	failedConfig := &domain.Config{
		Server: domain.Server{
			Port: 3000,
		},
		Key: "12345",
		Database: domain.Database{
			Host:     "test",
			Port:     1234,
			Username: "test",
			Password: "test",
			Name:     "test",
		},
		JWT: domain.JWT{
			Key:        "test",
			ExpireTime: 30,
		},
		Redis: domain.Redis{
			Host:     "test",
			Port:     1234,
			Password: "",
		},
	}

	db := &gorm.DB{}

	tests := []struct {
		name string
		args args
		mock func() AuthServiceInterface
		want error
	}{
		{
			name: "failed validating request",
			args: args{
				ctx: ctx,
				req: domain.AuthRequest{
					Email:    "",
					Password: "",
				},
			},
			mock: func() AuthServiceInterface {
				return &authServiceModule{
					cfg: failedConfig,
					db:  db,
				}
			},
			want: nil,
		},
		{
			name: "failed register",
			args: args{
				ctx: ctx,
				req: domain.AuthRequest{
					Email:    "test@mail.com",
					Password: "testtest",
				},
			},
			mock: func() AuthServiceInterface {
				authMock := NewMockAuthRepoInterface(ctrl)

				authMock.EXPECT().Register(ctx, gomock.Any()).Return(domain.User{}, errors.New("test"))

				return &authServiceModule{
					cfg:      failedConfig,
					db:       db,
					authRepo: authMock,
				}
			},
			want: errors.New("test"),
		},
		{
			name: "failed create inventory",
			args: args{
				ctx: ctx,
				req: domain.AuthRequest{
					Email:    "test@mail.com",
					Password: "testtest",
				},
			},
			mock: func() AuthServiceInterface {
				authMock := NewMockAuthRepoInterface(ctrl)
				authMock.EXPECT().Register(ctx, gomock.Any()).Return(domain.User{}, nil)

				inventoryMock := NewMockInventoryRepoInterface(ctrl)
				inventoryMock.EXPECT().Create(ctx, gomock.Any()).Return(errors.New("test"))

				return &authServiceModule{
					cfg:           failedConfig,
					db:            db,
					authRepo:      authMock,
					inventoryRepo: inventoryMock,
				}
			},
			want: errors.New("test"),
		},
		{
			name: "failed create profile",
			args: args{
				ctx: ctx,
				req: domain.AuthRequest{
					Email:    "test@mail.com",
					Password: "testtest",
				},
			},
			mock: func() AuthServiceInterface {
				authMock := NewMockAuthRepoInterface(ctrl)
				authMock.EXPECT().Register(ctx, gomock.Any()).Return(domain.User{}, nil)

				inventoryMock := NewMockInventoryRepoInterface(ctrl)
				inventoryMock.EXPECT().Create(ctx, gomock.Any()).Return(nil)

				profileMock := NewMockProfileRepoInterface(ctrl)
				profileMock.EXPECT().Create(ctx, gomock.Any(), gomock.Any()).Return(errors.New("test"))

				return &authServiceModule{
					cfg:           failedConfig,
					db:            db,
					authRepo:      authMock,
					inventoryRepo: inventoryMock,
					profileRepo:   profileMock,
				}
			},
			want: errors.New("test"),
		},
		{
			name: "success",
			args: args{
				ctx: ctx,
				req: domain.AuthRequest{
					Email:    "test@mail.com",
					Password: "testtest",
				},
			},
			mock: func() AuthServiceInterface {
				authMock := NewMockAuthRepoInterface(ctrl)
				authMock.EXPECT().Register(ctx, gomock.Any()).Return(domain.User{}, nil)

				inventoryMock := NewMockInventoryRepoInterface(ctrl)
				inventoryMock.EXPECT().Create(ctx, gomock.Any()).Return(nil)

				profileMock := NewMockProfileRepoInterface(ctrl)
				profileMock.EXPECT().Create(ctx, gomock.Any(), gomock.Any()).Return(nil)

				return &authServiceModule{
					cfg:           failedConfig,
					db:            db,
					authRepo:      authMock,
					inventoryRepo: inventoryMock,
					profileRepo:   profileMock,
				}
			},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := tt.mock()
			err := a.Register(tt.args.ctx, tt.args.req)
			if (err != nil) != (tt.want != nil) {
				t.Errorf("Register() error = %v, wantErr %v", err, tt.want)
				return
			}
		})
	}
}
