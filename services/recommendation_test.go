package services

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/zombozo12/tinder-dealls/domain"
	"gorm.io/gorm"
	"reflect"
	"testing"
)

func TestNewRecommendationService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		cfg         *domain.Config
		db          *gorm.DB
		redisRepo   RedisRepoInterface
		profileRepo ProfileRepoInterface
	}

	config := &domain.Config{}
	db := &gorm.DB{}
	redisMock := NewMockRedisRepoInterface(ctrl)
	profileMock := NewMockProfileRepoInterface(ctrl)

	tests := []struct {
		name    string
		args    args
		want    RecommendationServiceInterface
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				cfg:         config,
				db:          db,
				profileRepo: profileMock,
				redisRepo:   redisMock,
			},
			want: &recommendationServiceModule{
				cfg:         config,
				db:          db,
				profileRepo: profileMock,
				redisRepo:   redisMock,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewRecommendationService(tt.args.cfg, tt.args.db, tt.args.redisRepo, tt.args.profileRepo)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewRecommendationService() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got == nil {
				t.Errorf("NewRecommendationService() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRecommendationServiceModule_GetRecommendation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		ctx    context.Context
		userID int64
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

	ctx := context.WithValue(context.Background(), "requestid", "test")

	successArgs := args{
		ctx:    ctx,
		userID: 1,
	}

	tests := []struct {
		name    string
		args    args
		mock    func() *recommendationServiceModule
		want    []domain.Profile
		wantErr bool
	}{
		{
			name: "failed get profile",
			args: successArgs,
			mock: func() *recommendationServiceModule {
				profileMock := NewMockProfileRepoInterface(ctrl)
				profileMock.EXPECT().GetProfile(ctx, int64(1)).
					Return(nil, errors.New("test"))
				return &recommendationServiceModule{
					cfg:         config,
					db:          nil,
					profileRepo: profileMock,
					redisRepo:   nil,
				}
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "failed get frozen ids",
			args: successArgs,
			mock: func() *recommendationServiceModule {
				profileMock := NewMockProfileRepoInterface(ctrl)
				profileMock.EXPECT().GetProfile(ctx, int64(1)).
					Return(&domain.Profile{
						UserID: 1,
					}, nil)
				redisMock := NewMockRedisRepoInterface(ctrl)
				redisMock.EXPECT().Get(ctx, "frozen:1").
					Return("", errors.New("test"))
				return &recommendationServiceModule{
					cfg:         config,
					db:          nil,
					profileRepo: profileMock,
					redisRepo:   redisMock,
				}
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "failed unmarshal",
			args: successArgs,
			mock: func() *recommendationServiceModule {
				profileMock := NewMockProfileRepoInterface(ctrl)
				profileMock.EXPECT().GetProfile(ctx, int64(1)).
					Return(&domain.Profile{
						UserID: 1,
					}, nil)
				redisMock := NewMockRedisRepoInterface(ctrl)
				redisMock.EXPECT().Get(ctx, "frozen:1").
					Return("test", nil)
				return &recommendationServiceModule{
					cfg:         config,
					db:          nil,
					profileRepo: profileMock,
					redisRepo:   redisMock,
				}
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "failed get profile recommendation",
			args: successArgs,
			mock: func() *recommendationServiceModule {
				profileMock := NewMockProfileRepoInterface(ctrl)
				profileMock.EXPECT().GetProfile(ctx, int64(1)).
					Return(&domain.Profile{
						UserID: 1,
					}, nil)
				redisMock := NewMockRedisRepoInterface(ctrl)
				redisMock.EXPECT().Get(ctx, "frozen:1").
					Return("[1]", nil)
				profileMock.EXPECT().GetProfileRecommendation(ctx, "", []int64{1}).
					Return(nil, errors.New("test"))
				return &recommendationServiceModule{
					cfg:         config,
					db:          nil,
					profileRepo: profileMock,
					redisRepo:   redisMock,
				}
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "success",
			args: successArgs,
			mock: func() *recommendationServiceModule {
				profileMock := NewMockProfileRepoInterface(ctrl)
				profileMock.EXPECT().GetProfile(ctx, int64(1)).
					Return(&domain.Profile{
						UserID: 1,
					}, nil)
				redisMock := NewMockRedisRepoInterface(ctrl)
				redisMock.EXPECT().Get(ctx, "frozen:1").
					Return("[1]", nil)
				profileMock.EXPECT().GetProfileRecommendation(ctx, "", []int64{1}).
					Return([]domain.Profile{
						{
							UserID: 1,
						},
					}, nil)

				redisMock.EXPECT().Set(ctx, "frozen:1", gomock.Any(), 60*60*24).
					Return(nil)
				return &recommendationServiceModule{
					cfg:         config,
					db:          nil,
					profileRepo: profileMock,
					redisRepo:   redisMock,
				}
			},
			want: []domain.Profile{
				{
					UserID: 1,
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := tt.mock()
			got, err := m.GetRecommendation(tt.args.ctx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetRecommendation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetRecommendation() got = %v, want %v", got, tt.want)
			}
		})
	}
}
