package services

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/zombozo12/tinder-dealls/domain"
	"reflect"
	"testing"
)

func TestNewMatcherService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		cfg              *domain.Config
		profileRepo      ProfileRepoInterface
		inventoryRepo    InventoryRepoInterface
		matchedRepo      MatchedRepoInterface
		notificationRepo NotificationRepoInterface
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

	inventoryMock := NewMockInventoryRepoInterface(ctrl)
	profileMock := NewMockProfileRepoInterface(ctrl)
	matchedMock := NewMockMatchedRepoInterface(ctrl)
	notificationMock := NewMockNotificationRepoInterface(ctrl)

	tests := []struct {
		name    string
		args    args
		want    MatcherServiceInterface
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				cfg:              config,
				profileRepo:      profileMock,
				inventoryRepo:    inventoryMock,
				matchedRepo:      matchedMock,
				notificationRepo: notificationMock,
			},
			want: &matcherServiceModule{
				cfg:              config,
				profileRepo:      profileMock,
				inventoryRepo:    inventoryMock,
				matchedRepo:      matchedMock,
				notificationRepo: notificationMock,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewMatcherService(tt.args.cfg, tt.args.profileRepo, tt.args.inventoryRepo, tt.args.matchedRepo, tt.args.notificationRepo)
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

func Test_matcherServiceModule_Like(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

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

	type args struct {
		ctx          context.Context
		userID       int64
		targetUserID int64
	}

	failedArgs := args{
		ctx:          ctx,
		userID:       1,
		targetUserID: 1,
	}

	successArgs := args{
		ctx:          ctx,
		userID:       1,
		targetUserID: 2,
	}

	tests := []struct {
		name    string
		mock    func() *matcherServiceModule
		args    args
		wantErr bool
	}{
		{
			name: "failed target user id is same with user id",
			mock: func() *matcherServiceModule {
				return &matcherServiceModule{
					cfg:              config,
					profileRepo:      nil,
					inventoryRepo:    nil,
					matchedRepo:      nil,
					notificationRepo: nil,
				}
			},
			args:    failedArgs,
			wantErr: true,
		},
		{
			name: "failed get inventory by user id",
			mock: func() *matcherServiceModule {
				inventoryMock := NewMockInventoryRepoInterface(ctrl)
				inventoryMock.EXPECT().GetByUserId(ctx, int64(1)).Return(nil, errors.New("test"))

				return &matcherServiceModule{
					cfg:              config,
					profileRepo:      nil,
					inventoryRepo:    inventoryMock,
					matchedRepo:      nil,
					notificationRepo: nil,
				}
			},
			args:    successArgs,
			wantErr: true,
		},
		{
			name: "failed insufficient swipes",
			mock: func() *matcherServiceModule {
				inventoryMock := NewMockInventoryRepoInterface(ctrl)
				inventoryMock.EXPECT().GetByUserId(ctx, int64(1)).Return(&domain.Inventory{
					Swipes: 0,
				}, nil)

				return &matcherServiceModule{
					cfg:              config,
					profileRepo:      nil,
					inventoryRepo:    inventoryMock,
					matchedRepo:      nil,
					notificationRepo: nil,
				}
			},
			args:    successArgs,
			wantErr: true,
		},
		{
			name: "failed insufficient likes",
			mock: func() *matcherServiceModule {
				inventoryMock := NewMockInventoryRepoInterface(ctrl)
				inventoryMock.EXPECT().GetByUserId(ctx, int64(1)).Return(&domain.Inventory{
					Swipes: 1,
					Likes:  0,
				}, nil)

				return &matcherServiceModule{
					cfg:              config,
					profileRepo:      nil,
					inventoryRepo:    inventoryMock,
					matchedRepo:      nil,
					notificationRepo: nil,
				}
			},
			args:    successArgs,
			wantErr: true,
		},
		{
			name: "failed get profile",
			mock: func() *matcherServiceModule {
				inventoryMock := NewMockInventoryRepoInterface(ctrl)
				inventoryMock.EXPECT().GetByUserId(ctx, int64(1)).Return(&domain.Inventory{
					Swipes: 1,
					Likes:  1,
				}, nil)

				profileMock := NewMockProfileRepoInterface(ctrl)
				profileMock.EXPECT().GetProfile(ctx, int64(1)).Return(nil, errors.New("test"))

				return &matcherServiceModule{
					cfg:              config,
					profileRepo:      profileMock,
					inventoryRepo:    inventoryMock,
					matchedRepo:      nil,
					notificationRepo: nil,
				}
			},
			args:    successArgs,
			wantErr: true,
		},
		{
			name: "failed profile not found",
			mock: func() *matcherServiceModule {
				inventoryMock := NewMockInventoryRepoInterface(ctrl)
				inventoryMock.EXPECT().GetByUserId(ctx, int64(1)).Return(&domain.Inventory{
					Swipes: 1,
					Likes:  1,
				}, nil)

				profileMock := NewMockProfileRepoInterface(ctrl)
				profileMock.EXPECT().GetProfile(ctx, int64(1)).Return(nil, nil)

				return &matcherServiceModule{
					cfg:              config,
					profileRepo:      profileMock,
					inventoryRepo:    inventoryMock,
					matchedRepo:      nil,
					notificationRepo: nil,
				}
			},
			args:    successArgs,
			wantErr: true,
		},
		{
			name: "failed is exists",
			mock: func() *matcherServiceModule {
				inventoryMock := NewMockInventoryRepoInterface(ctrl)
				inventoryMock.EXPECT().GetByUserId(ctx, int64(1)).Return(&domain.Inventory{
					Swipes: 1,
					Likes:  1,
				}, nil)

				profileMock := NewMockProfileRepoInterface(ctrl)
				profileMock.EXPECT().GetProfile(ctx, int64(1)).Return(&domain.Profile{
					ID: 1,
				}, nil)

				matchedMock := NewMockMatchedRepoInterface(ctrl)
				matchedMock.EXPECT().IsExists(ctx, domain.MatchRequest{
					UserID:       int64(1),
					TargetUserID: int64(2),
				}).Return(false, errors.New("test"))

				return &matcherServiceModule{
					cfg:              config,
					profileRepo:      profileMock,
					inventoryRepo:    inventoryMock,
					matchedRepo:      matchedMock,
					notificationRepo: nil,
				}
			},
			args:    successArgs,
			wantErr: true,
		},
		{
			name: "match is exists",
			mock: func() *matcherServiceModule {
				inventoryMock := NewMockInventoryRepoInterface(ctrl)
				inventoryMock.EXPECT().GetByUserId(ctx, int64(1)).Return(&domain.Inventory{
					Swipes: 1,
					Likes:  1,
				}, nil)

				profileMock := NewMockProfileRepoInterface(ctrl)
				profileMock.EXPECT().GetProfile(ctx, int64(1)).Return(&domain.Profile{
					ID: 1,
				}, nil)

				matchedMock := NewMockMatchedRepoInterface(ctrl)
				matchedMock.EXPECT().IsExists(ctx, domain.MatchRequest{
					UserID:       int64(1),
					TargetUserID: int64(2),
				}).Return(true, nil)

				return &matcherServiceModule{
					cfg:              config,
					profileRepo:      profileMock,
					inventoryRepo:    inventoryMock,
					matchedRepo:      matchedMock,
					notificationRepo: nil,
				}
			},
			args:    successArgs,
			wantErr: true,
		},
		{
			name: "failed create match",
			mock: func() *matcherServiceModule {
				inventoryMock := NewMockInventoryRepoInterface(ctrl)
				inventoryMock.EXPECT().GetByUserId(ctx, int64(1)).Return(&domain.Inventory{
					Swipes: 1,
					Likes:  1,
				}, nil)

				profileMock := NewMockProfileRepoInterface(ctrl)
				profileMock.EXPECT().GetProfile(ctx, int64(1)).Return(&domain.Profile{
					ID: 1,
				}, nil)

				matchedMock := NewMockMatchedRepoInterface(ctrl)
				matchedMock.EXPECT().IsExists(ctx, domain.MatchRequest{
					UserID:       int64(1),
					TargetUserID: int64(2),
				}).Return(false, nil)
				matchedMock.EXPECT().Create(ctx, domain.MatchRequest{
					UserID:       int64(1),
					TargetUserID: int64(2),
				}).Return(errors.New("test"))

				return &matcherServiceModule{
					cfg:              config,
					profileRepo:      profileMock,
					inventoryRepo:    inventoryMock,
					matchedRepo:      matchedMock,
					notificationRepo: nil,
				}
			},
			args:    successArgs,
			wantErr: true,
		},
		{
			name: "failed get is matched",
			mock: func() *matcherServiceModule {
				inventoryMock := NewMockInventoryRepoInterface(ctrl)
				inventoryMock.EXPECT().GetByUserId(ctx, int64(1)).Return(&domain.Inventory{
					Swipes: 1,
					Likes:  1,
				}, nil)

				profileMock := NewMockProfileRepoInterface(ctrl)
				profileMock.EXPECT().GetProfile(ctx, int64(1)).Return(&domain.Profile{
					ID: 1,
				}, nil)

				matchedMock := NewMockMatchedRepoInterface(ctrl)
				matchedMock.EXPECT().IsExists(ctx, domain.MatchRequest{
					UserID:       int64(1),
					TargetUserID: int64(2),
				}).Return(false, nil)
				matchedMock.EXPECT().Create(ctx, domain.MatchRequest{
					UserID:       int64(1),
					TargetUserID: int64(2),
				}).Return(nil)
				matchedMock.EXPECT().IsMatched(ctx, domain.MatchRequest{
					UserID:       int64(2),
					TargetUserID: int64(1),
				}).Return(false, errors.New("test"))

				return &matcherServiceModule{
					cfg:              config,
					profileRepo:      profileMock,
					inventoryRepo:    inventoryMock,
					matchedRepo:      matchedMock,
					notificationRepo: nil,
				}
			},
			args:    successArgs,
			wantErr: true,
		},
		{
			name:    "is matched - failed create notification",
			args:    successArgs,
			wantErr: true,
			mock: func() *matcherServiceModule {
				inventoryMock := NewMockInventoryRepoInterface(ctrl)
				inventoryMock.EXPECT().GetByUserId(ctx, int64(1)).Return(&domain.Inventory{
					Swipes: 1,
					Likes:  1,
				}, nil)

				profileMock := NewMockProfileRepoInterface(ctrl)
				profileMock.EXPECT().GetProfile(ctx, int64(1)).Return(&domain.Profile{
					ID: 1,
				}, nil)

				matchedMock := NewMockMatchedRepoInterface(ctrl)
				matchedMock.EXPECT().IsExists(ctx, domain.MatchRequest{
					UserID:       int64(1),
					TargetUserID: int64(2),
				}).Return(false, nil)
				matchedMock.EXPECT().Create(ctx, domain.MatchRequest{
					UserID:       int64(1),
					TargetUserID: int64(2),
				}).Return(nil)
				matchedMock.EXPECT().IsMatched(ctx, domain.MatchRequest{
					UserID:       int64(2),
					TargetUserID: int64(1),
				}).Return(true, nil)

				notificationMock := NewMockNotificationRepoInterface(ctrl)
				notificationMock.EXPECT().Create(ctx, domain.NotificationRequest{
					UserID:  int64(2),
					Message: "You have a new match with ",
				}).Return(errors.New("test"))

				return &matcherServiceModule{
					cfg:              config,
					profileRepo:      profileMock,
					inventoryRepo:    inventoryMock,
					matchedRepo:      matchedMock,
					notificationRepo: notificationMock,
				}
			},
		},
		{
			name:    "is matched - success create notification",
			args:    successArgs,
			wantErr: false,
			mock: func() *matcherServiceModule {
				inventoryMock := NewMockInventoryRepoInterface(ctrl)
				inventoryMock.EXPECT().GetByUserId(ctx, int64(1)).Return(&domain.Inventory{
					Swipes: 1,
					Likes:  1,
				}, nil)

				profileMock := NewMockProfileRepoInterface(ctrl)
				profileMock.EXPECT().GetProfile(ctx, int64(1)).Return(&domain.Profile{
					ID: 1,
				}, nil)

				matchedMock := NewMockMatchedRepoInterface(ctrl)
				matchedMock.EXPECT().IsExists(ctx, domain.MatchRequest{
					UserID:       int64(1),
					TargetUserID: int64(2),
				}).Return(false, nil)
				matchedMock.EXPECT().Create(ctx, domain.MatchRequest{
					UserID:       int64(1),
					TargetUserID: int64(2),
				}).Return(nil)
				matchedMock.EXPECT().IsMatched(ctx, domain.MatchRequest{
					UserID:       int64(2),
					TargetUserID: int64(1),
				}).Return(true, nil)

				notificationMock := NewMockNotificationRepoInterface(ctrl)
				notificationMock.EXPECT().Create(ctx, domain.NotificationRequest{
					UserID:  int64(2),
					Message: "You have a new match with ",
				}).Return(nil)

				return &matcherServiceModule{
					cfg:              config,
					profileRepo:      profileMock,
					inventoryRepo:    inventoryMock,
					matchedRepo:      matchedMock,
					notificationRepo: notificationMock,
				}
			},
		},
		{
			name:    "failed update swipes",
			args:    successArgs,
			wantErr: true,
			mock: func() *matcherServiceModule {
				inventoryMock := NewMockInventoryRepoInterface(ctrl)
				inventoryMock.EXPECT().GetByUserId(ctx, int64(1)).Return(&domain.Inventory{
					ID:     1,
					Swipes: 1,
					Likes:  1,
				}, nil)
				profileMock := NewMockProfileRepoInterface(ctrl)
				profileMock.EXPECT().GetProfile(ctx, int64(1)).Return(&domain.Profile{
					ID: 1,
				}, nil)

				matchedMock := NewMockMatchedRepoInterface(ctrl)
				matchedMock.EXPECT().IsExists(ctx, domain.MatchRequest{
					UserID:       int64(1),
					TargetUserID: int64(2),
				}).Return(false, nil)
				matchedMock.EXPECT().Create(ctx, domain.MatchRequest{
					UserID:       int64(1),
					TargetUserID: int64(2),
				}).Return(nil)
				matchedMock.EXPECT().IsMatched(ctx, domain.MatchRequest{
					UserID:       int64(2),
					TargetUserID: int64(1),
				}).Return(false, nil)

				inventoryMock.EXPECT().UpdateSwipes(ctx, int64(1), 0).
					Return(errors.New("test"))

				return &matcherServiceModule{
					cfg:              config,
					profileRepo:      profileMock,
					inventoryRepo:    inventoryMock,
					matchedRepo:      matchedMock,
					notificationRepo: nil,
				}
			},
		},
		{
			name:    "failed update swipes",
			args:    successArgs,
			wantErr: true,
			mock: func() *matcherServiceModule {
				inventoryMock := NewMockInventoryRepoInterface(ctrl)
				inventoryMock.EXPECT().GetByUserId(ctx, int64(1)).Return(&domain.Inventory{
					ID:     1,
					Swipes: 1,
					Likes:  1,
				}, nil)
				profileMock := NewMockProfileRepoInterface(ctrl)
				profileMock.EXPECT().GetProfile(ctx, int64(1)).Return(&domain.Profile{
					ID: 1,
				}, nil)

				matchedMock := NewMockMatchedRepoInterface(ctrl)
				matchedMock.EXPECT().IsExists(ctx, domain.MatchRequest{
					UserID:       int64(1),
					TargetUserID: int64(2),
				}).Return(false, nil)
				matchedMock.EXPECT().Create(ctx, domain.MatchRequest{
					UserID:       int64(1),
					TargetUserID: int64(2),
				}).Return(nil)
				matchedMock.EXPECT().IsMatched(ctx, domain.MatchRequest{
					UserID:       int64(2),
					TargetUserID: int64(1),
				}).Return(false, nil)

				inventoryMock.EXPECT().UpdateSwipes(ctx, int64(1), 0).
					Return(errors.New("test"))

				return &matcherServiceModule{
					cfg:              config,
					profileRepo:      profileMock,
					inventoryRepo:    inventoryMock,
					matchedRepo:      matchedMock,
					notificationRepo: nil,
				}
			},
		},
		{
			name:    "failed update likes",
			args:    successArgs,
			wantErr: true,
			mock: func() *matcherServiceModule {
				inventoryMock := NewMockInventoryRepoInterface(ctrl)
				inventoryMock.EXPECT().GetByUserId(ctx, int64(1)).Return(&domain.Inventory{
					ID:     1,
					Swipes: 1,
					Likes:  1,
				}, nil)
				profileMock := NewMockProfileRepoInterface(ctrl)
				profileMock.EXPECT().GetProfile(ctx, int64(1)).Return(&domain.Profile{
					ID: 1,
				}, nil)

				matchedMock := NewMockMatchedRepoInterface(ctrl)
				matchedMock.EXPECT().IsExists(ctx, domain.MatchRequest{
					UserID:       int64(1),
					TargetUserID: int64(2),
				}).Return(false, nil)
				matchedMock.EXPECT().Create(ctx, domain.MatchRequest{
					UserID:       int64(1),
					TargetUserID: int64(2),
				}).Return(nil)
				matchedMock.EXPECT().IsMatched(ctx, domain.MatchRequest{
					UserID:       int64(2),
					TargetUserID: int64(1),
				}).Return(false, nil)

				inventoryMock.EXPECT().UpdateSwipes(ctx, int64(1), 0).
					Return(nil)

				inventoryMock.EXPECT().UpdateLikes(ctx, int64(1), 0).
					Return(errors.New("test"))

				return &matcherServiceModule{
					cfg:              config,
					profileRepo:      profileMock,
					inventoryRepo:    inventoryMock,
					matchedRepo:      matchedMock,
					notificationRepo: nil,
				}
			},
		},
		{
			name:    "success",
			args:    successArgs,
			wantErr: false,
			mock: func() *matcherServiceModule {
				inventoryMock := NewMockInventoryRepoInterface(ctrl)
				inventoryMock.EXPECT().GetByUserId(ctx, int64(1)).Return(&domain.Inventory{
					ID:     1,
					Swipes: 1,
					Likes:  1,
				}, nil)
				profileMock := NewMockProfileRepoInterface(ctrl)
				profileMock.EXPECT().GetProfile(ctx, int64(1)).Return(&domain.Profile{
					ID: 1,
				}, nil)

				matchedMock := NewMockMatchedRepoInterface(ctrl)
				matchedMock.EXPECT().IsExists(ctx, domain.MatchRequest{
					UserID:       int64(1),
					TargetUserID: int64(2),
				}).Return(false, nil)
				matchedMock.EXPECT().Create(ctx, domain.MatchRequest{
					UserID:       int64(1),
					TargetUserID: int64(2),
				}).Return(nil)
				matchedMock.EXPECT().IsMatched(ctx, domain.MatchRequest{
					UserID:       int64(2),
					TargetUserID: int64(1),
				}).Return(false, nil)

				inventoryMock.EXPECT().UpdateSwipes(ctx, int64(1), 0).
					Return(nil)

				inventoryMock.EXPECT().UpdateLikes(ctx, int64(1), 0).
					Return(nil)

				return &matcherServiceModule{
					cfg:              config,
					profileRepo:      profileMock,
					inventoryRepo:    inventoryMock,
					matchedRepo:      matchedMock,
					notificationRepo: nil,
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := tt.mock()
			if err := m.Like(tt.args.ctx, tt.args.userID, tt.args.targetUserID); (err != nil) != tt.wantErr {
				t.Errorf("matcherServiceModule.Like() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_matcherServiceModule_SuperLike(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

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

	type args struct {
		ctx          context.Context
		userID       int64
		targetUserID int64
	}

	failedArgs := args{
		ctx:          ctx,
		userID:       1,
		targetUserID: 1,
	}

	successArgs := args{
		ctx:          ctx,
		userID:       1,
		targetUserID: 2,
	}

	tests := []struct {
		name    string
		mock    func() *matcherServiceModule
		args    args
		wantErr bool
	}{
		{
			name: "failed target user id is same with user id",
			mock: func() *matcherServiceModule {
				return &matcherServiceModule{
					cfg:              config,
					profileRepo:      nil,
					inventoryRepo:    nil,
					matchedRepo:      nil,
					notificationRepo: nil,
				}
			},
			args:    failedArgs,
			wantErr: true,
		},
		{
			name: "failed get inventory by user id",
			mock: func() *matcherServiceModule {
				inventoryMock := NewMockInventoryRepoInterface(ctrl)
				inventoryMock.EXPECT().GetByUserId(ctx, int64(1)).Return(nil, errors.New("test"))

				return &matcherServiceModule{
					cfg:              config,
					profileRepo:      nil,
					inventoryRepo:    inventoryMock,
					matchedRepo:      nil,
					notificationRepo: nil,
				}
			},
			args:    successArgs,
			wantErr: true,
		},
		{
			name: "failed insufficient swipes",
			mock: func() *matcherServiceModule {
				inventoryMock := NewMockInventoryRepoInterface(ctrl)
				inventoryMock.EXPECT().GetByUserId(ctx, int64(1)).Return(&domain.Inventory{
					Swipes: 0,
				}, nil)

				return &matcherServiceModule{
					cfg:              config,
					profileRepo:      nil,
					inventoryRepo:    inventoryMock,
					matchedRepo:      nil,
					notificationRepo: nil,
				}
			},
			args:    successArgs,
			wantErr: true,
		},
		{
			name: "failed insufficient superlikes",
			mock: func() *matcherServiceModule {
				inventoryMock := NewMockInventoryRepoInterface(ctrl)
				inventoryMock.EXPECT().GetByUserId(ctx, int64(1)).Return(&domain.Inventory{
					Swipes:     1,
					SuperLikes: 0,
				}, nil)

				return &matcherServiceModule{
					cfg:              config,
					profileRepo:      nil,
					inventoryRepo:    inventoryMock,
					matchedRepo:      nil,
					notificationRepo: nil,
				}
			},
			args:    successArgs,
			wantErr: true,
		},
		{
			name: "failed get profile",
			mock: func() *matcherServiceModule {
				inventoryMock := NewMockInventoryRepoInterface(ctrl)
				inventoryMock.EXPECT().GetByUserId(ctx, int64(1)).Return(&domain.Inventory{
					Swipes:     1,
					SuperLikes: 1,
				}, nil)

				profileMock := NewMockProfileRepoInterface(ctrl)
				profileMock.EXPECT().GetProfile(ctx, int64(1)).Return(nil, errors.New("test"))

				return &matcherServiceModule{
					cfg:              config,
					profileRepo:      profileMock,
					inventoryRepo:    inventoryMock,
					matchedRepo:      nil,
					notificationRepo: nil,
				}
			},
			args:    successArgs,
			wantErr: true,
		},
		{
			name: "failed profile not found",
			mock: func() *matcherServiceModule {
				inventoryMock := NewMockInventoryRepoInterface(ctrl)
				inventoryMock.EXPECT().GetByUserId(ctx, int64(1)).Return(&domain.Inventory{
					Swipes:     1,
					SuperLikes: 1,
				}, nil)

				profileMock := NewMockProfileRepoInterface(ctrl)
				profileMock.EXPECT().GetProfile(ctx, int64(1)).Return(nil, nil)

				return &matcherServiceModule{
					cfg:              config,
					profileRepo:      profileMock,
					inventoryRepo:    inventoryMock,
					matchedRepo:      nil,
					notificationRepo: nil,
				}
			},
			args:    successArgs,
			wantErr: true,
		},
		{
			name: "failed is exists",
			mock: func() *matcherServiceModule {
				inventoryMock := NewMockInventoryRepoInterface(ctrl)
				inventoryMock.EXPECT().GetByUserId(ctx, int64(1)).Return(&domain.Inventory{
					Swipes:     1,
					SuperLikes: 1,
				}, nil)

				profileMock := NewMockProfileRepoInterface(ctrl)
				profileMock.EXPECT().GetProfile(ctx, int64(1)).Return(&domain.Profile{
					ID: 1,
				}, nil)

				matchedMock := NewMockMatchedRepoInterface(ctrl)
				matchedMock.EXPECT().IsExists(ctx, domain.MatchRequest{
					UserID:       int64(1),
					TargetUserID: int64(2),
				}).Return(false, errors.New("test"))

				return &matcherServiceModule{
					cfg:              config,
					profileRepo:      profileMock,
					inventoryRepo:    inventoryMock,
					matchedRepo:      matchedMock,
					notificationRepo: nil,
				}
			},
			args:    successArgs,
			wantErr: true,
		},
		{
			name: "match is exists",
			mock: func() *matcherServiceModule {
				inventoryMock := NewMockInventoryRepoInterface(ctrl)
				inventoryMock.EXPECT().GetByUserId(ctx, int64(1)).Return(&domain.Inventory{
					Swipes:     1,
					SuperLikes: 1,
				}, nil)

				profileMock := NewMockProfileRepoInterface(ctrl)
				profileMock.EXPECT().GetProfile(ctx, int64(1)).Return(&domain.Profile{
					ID: 1,
				}, nil)

				matchedMock := NewMockMatchedRepoInterface(ctrl)
				matchedMock.EXPECT().IsExists(ctx, domain.MatchRequest{
					UserID:       int64(1),
					TargetUserID: int64(2),
				}).Return(true, nil)

				return &matcherServiceModule{
					cfg:              config,
					profileRepo:      profileMock,
					inventoryRepo:    inventoryMock,
					matchedRepo:      matchedMock,
					notificationRepo: nil,
				}
			},
			args:    successArgs,
			wantErr: true,
		},
		{
			name: "failed create match",
			mock: func() *matcherServiceModule {
				inventoryMock := NewMockInventoryRepoInterface(ctrl)
				inventoryMock.EXPECT().GetByUserId(ctx, int64(1)).Return(&domain.Inventory{
					Swipes:     1,
					SuperLikes: 1,
				}, nil)

				profileMock := NewMockProfileRepoInterface(ctrl)
				profileMock.EXPECT().GetProfile(ctx, int64(1)).Return(&domain.Profile{
					ID: 1,
				}, nil)

				matchedMock := NewMockMatchedRepoInterface(ctrl)
				matchedMock.EXPECT().IsExists(ctx, domain.MatchRequest{
					UserID:       int64(1),
					TargetUserID: int64(2),
				}).Return(false, nil)
				matchedMock.EXPECT().Create(ctx, domain.MatchRequest{
					UserID:       int64(1),
					TargetUserID: int64(2),
				}).Return(errors.New("test"))

				return &matcherServiceModule{
					cfg:              config,
					profileRepo:      profileMock,
					inventoryRepo:    inventoryMock,
					matchedRepo:      matchedMock,
					notificationRepo: nil,
				}
			},
			args:    successArgs,
			wantErr: true,
		},
		{
			name: "failed get is matched",
			mock: func() *matcherServiceModule {
				inventoryMock := NewMockInventoryRepoInterface(ctrl)
				inventoryMock.EXPECT().GetByUserId(ctx, int64(1)).Return(&domain.Inventory{
					Swipes:     1,
					SuperLikes: 1,
				}, nil)

				profileMock := NewMockProfileRepoInterface(ctrl)
				profileMock.EXPECT().GetProfile(ctx, int64(1)).Return(&domain.Profile{
					ID: 1,
				}, nil)

				matchedMock := NewMockMatchedRepoInterface(ctrl)
				matchedMock.EXPECT().IsExists(ctx, domain.MatchRequest{
					UserID:       int64(1),
					TargetUserID: int64(2),
				}).Return(false, nil)
				matchedMock.EXPECT().Create(ctx, domain.MatchRequest{
					UserID:       int64(1),
					TargetUserID: int64(2),
				}).Return(nil)
				matchedMock.EXPECT().IsMatched(ctx, domain.MatchRequest{
					UserID:       int64(2),
					TargetUserID: int64(1),
				}).Return(false, errors.New("test"))

				return &matcherServiceModule{
					cfg:              config,
					profileRepo:      profileMock,
					inventoryRepo:    inventoryMock,
					matchedRepo:      matchedMock,
					notificationRepo: nil,
				}
			},
			args:    successArgs,
			wantErr: true,
		},
		{
			name:    "is matched - failed create notification",
			args:    successArgs,
			wantErr: true,
			mock: func() *matcherServiceModule {
				inventoryMock := NewMockInventoryRepoInterface(ctrl)
				inventoryMock.EXPECT().GetByUserId(ctx, int64(1)).Return(&domain.Inventory{
					Swipes:     1,
					SuperLikes: 1,
				}, nil)

				profileMock := NewMockProfileRepoInterface(ctrl)
				profileMock.EXPECT().GetProfile(ctx, int64(1)).Return(&domain.Profile{
					ID: 1,
				}, nil)

				matchedMock := NewMockMatchedRepoInterface(ctrl)
				matchedMock.EXPECT().IsExists(ctx, domain.MatchRequest{
					UserID:       int64(1),
					TargetUserID: int64(2),
				}).Return(false, nil)
				matchedMock.EXPECT().Create(ctx, domain.MatchRequest{
					UserID:       int64(1),
					TargetUserID: int64(2),
				}).Return(nil)
				matchedMock.EXPECT().IsMatched(ctx, domain.MatchRequest{
					UserID:       int64(2),
					TargetUserID: int64(1),
				}).Return(true, nil)

				notificationMock := NewMockNotificationRepoInterface(ctrl)
				notificationMock.EXPECT().Create(ctx, domain.NotificationRequest{
					UserID:  int64(2),
					Message: "You have a new match with ",
				}).Return(errors.New("test"))

				return &matcherServiceModule{
					cfg:              config,
					profileRepo:      profileMock,
					inventoryRepo:    inventoryMock,
					matchedRepo:      matchedMock,
					notificationRepo: notificationMock,
				}
			},
		},
		{
			name:    "is matched - success create notification",
			args:    successArgs,
			wantErr: false,
			mock: func() *matcherServiceModule {
				inventoryMock := NewMockInventoryRepoInterface(ctrl)
				inventoryMock.EXPECT().GetByUserId(ctx, int64(1)).Return(&domain.Inventory{
					Swipes:     1,
					SuperLikes: 1,
				}, nil)

				profileMock := NewMockProfileRepoInterface(ctrl)
				profileMock.EXPECT().GetProfile(ctx, int64(1)).Return(&domain.Profile{
					ID: 1,
				}, nil)

				matchedMock := NewMockMatchedRepoInterface(ctrl)
				matchedMock.EXPECT().IsExists(ctx, domain.MatchRequest{
					UserID:       int64(1),
					TargetUserID: int64(2),
				}).Return(false, nil)
				matchedMock.EXPECT().Create(ctx, domain.MatchRequest{
					UserID:       int64(1),
					TargetUserID: int64(2),
				}).Return(nil)
				matchedMock.EXPECT().IsMatched(ctx, domain.MatchRequest{
					UserID:       int64(2),
					TargetUserID: int64(1),
				}).Return(true, nil)

				notificationMock := NewMockNotificationRepoInterface(ctrl)
				notificationMock.EXPECT().Create(ctx, domain.NotificationRequest{
					UserID:  int64(2),
					Message: "You have a new match with ",
				}).Return(nil)

				return &matcherServiceModule{
					cfg:              config,
					profileRepo:      profileMock,
					inventoryRepo:    inventoryMock,
					matchedRepo:      matchedMock,
					notificationRepo: notificationMock,
				}
			},
		},
		{
			name:    "failed notification super liked",
			args:    successArgs,
			wantErr: true,
			mock: func() *matcherServiceModule {
				inventoryMock := NewMockInventoryRepoInterface(ctrl)
				inventoryMock.EXPECT().GetByUserId(ctx, int64(1)).Return(&domain.Inventory{
					ID:         1,
					Swipes:     1,
					SuperLikes: 1,
				}, nil)
				profileMock := NewMockProfileRepoInterface(ctrl)
				profileMock.EXPECT().GetProfile(ctx, int64(1)).Return(&domain.Profile{
					ID: 1,
				}, nil)

				matchedMock := NewMockMatchedRepoInterface(ctrl)
				matchedMock.EXPECT().IsExists(ctx, domain.MatchRequest{
					UserID:       int64(1),
					TargetUserID: int64(2),
				}).Return(false, nil)
				matchedMock.EXPECT().Create(ctx, domain.MatchRequest{
					UserID:       int64(1),
					TargetUserID: int64(2),
				}).Return(nil)
				matchedMock.EXPECT().IsMatched(ctx, domain.MatchRequest{
					UserID:       int64(2),
					TargetUserID: int64(1),
				}).Return(false, nil)

				notificationMock := NewMockNotificationRepoInterface(ctrl)
				notificationMock.EXPECT().Create(ctx, domain.NotificationRequest{
					UserID:  int64(2),
					Message: " super liked you",
				}).Return(errors.New("test"))

				return &matcherServiceModule{
					cfg:              config,
					profileRepo:      profileMock,
					inventoryRepo:    inventoryMock,
					matchedRepo:      matchedMock,
					notificationRepo: notificationMock,
				}
			},
		},
		{
			name:    "failed update swipes",
			args:    successArgs,
			wantErr: true,
			mock: func() *matcherServiceModule {
				inventoryMock := NewMockInventoryRepoInterface(ctrl)
				inventoryMock.EXPECT().GetByUserId(ctx, int64(1)).Return(&domain.Inventory{
					ID:         1,
					Swipes:     1,
					SuperLikes: 1,
				}, nil)
				profileMock := NewMockProfileRepoInterface(ctrl)
				profileMock.EXPECT().GetProfile(ctx, int64(1)).Return(&domain.Profile{
					ID: 1,
				}, nil)

				matchedMock := NewMockMatchedRepoInterface(ctrl)
				matchedMock.EXPECT().IsExists(ctx, domain.MatchRequest{
					UserID:       int64(1),
					TargetUserID: int64(2),
				}).Return(false, nil)
				matchedMock.EXPECT().Create(ctx, domain.MatchRequest{
					UserID:       int64(1),
					TargetUserID: int64(2),
				}).Return(nil)
				matchedMock.EXPECT().IsMatched(ctx, domain.MatchRequest{
					UserID:       int64(2),
					TargetUserID: int64(1),
				}).Return(false, nil)

				notificationMock := NewMockNotificationRepoInterface(ctrl)
				notificationMock.EXPECT().Create(ctx, domain.NotificationRequest{
					UserID:  int64(2),
					Message: " super liked you",
				}).Return(nil)

				inventoryMock.EXPECT().UpdateSwipes(ctx, int64(1), 0).
					Return(errors.New("test"))

				return &matcherServiceModule{
					cfg:              config,
					profileRepo:      profileMock,
					inventoryRepo:    inventoryMock,
					matchedRepo:      matchedMock,
					notificationRepo: notificationMock,
				}
			},
		},
		{
			name:    "failed update super likes",
			args:    successArgs,
			wantErr: true,
			mock: func() *matcherServiceModule {
				inventoryMock := NewMockInventoryRepoInterface(ctrl)
				inventoryMock.EXPECT().GetByUserId(ctx, int64(1)).Return(&domain.Inventory{
					ID:         1,
					Swipes:     1,
					SuperLikes: 1,
				}, nil)
				profileMock := NewMockProfileRepoInterface(ctrl)
				profileMock.EXPECT().GetProfile(ctx, int64(1)).Return(&domain.Profile{
					ID: 1,
				}, nil)

				matchedMock := NewMockMatchedRepoInterface(ctrl)
				matchedMock.EXPECT().IsExists(ctx, domain.MatchRequest{
					UserID:       int64(1),
					TargetUserID: int64(2),
				}).Return(false, nil)
				matchedMock.EXPECT().Create(ctx, domain.MatchRequest{
					UserID:       int64(1),
					TargetUserID: int64(2),
				}).Return(nil)
				matchedMock.EXPECT().IsMatched(ctx, domain.MatchRequest{
					UserID:       int64(2),
					TargetUserID: int64(1),
				}).Return(false, nil)

				notificationMock := NewMockNotificationRepoInterface(ctrl)
				notificationMock.EXPECT().Create(ctx, domain.NotificationRequest{
					UserID:  int64(2),
					Message: " super liked you",
				}).Return(nil)

				inventoryMock.EXPECT().UpdateSwipes(ctx, int64(1), 0).
					Return(nil)

				inventoryMock.EXPECT().UpdateSuperLikes(ctx, int64(1), 0).
					Return(errors.New("test"))

				return &matcherServiceModule{
					cfg:              config,
					profileRepo:      profileMock,
					inventoryRepo:    inventoryMock,
					matchedRepo:      matchedMock,
					notificationRepo: notificationMock,
				}
			},
		},
		{
			name:    "success",
			args:    successArgs,
			wantErr: false,
			mock: func() *matcherServiceModule {
				inventoryMock := NewMockInventoryRepoInterface(ctrl)
				inventoryMock.EXPECT().GetByUserId(ctx, int64(1)).Return(&domain.Inventory{
					ID:         1,
					Swipes:     1,
					SuperLikes: 1,
				}, nil)
				profileMock := NewMockProfileRepoInterface(ctrl)
				profileMock.EXPECT().GetProfile(ctx, int64(1)).Return(&domain.Profile{
					ID: 1,
				}, nil)

				matchedMock := NewMockMatchedRepoInterface(ctrl)
				matchedMock.EXPECT().IsExists(ctx, domain.MatchRequest{
					UserID:       int64(1),
					TargetUserID: int64(2),
				}).Return(false, nil)
				matchedMock.EXPECT().Create(ctx, domain.MatchRequest{
					UserID:       int64(1),
					TargetUserID: int64(2),
				}).Return(nil)
				matchedMock.EXPECT().IsMatched(ctx, domain.MatchRequest{
					UserID:       int64(2),
					TargetUserID: int64(1),
				}).Return(false, nil)

				notificationMock := NewMockNotificationRepoInterface(ctrl)
				notificationMock.EXPECT().Create(ctx, domain.NotificationRequest{
					UserID:  int64(2),
					Message: " super liked you",
				}).Return(nil)

				inventoryMock.EXPECT().UpdateSwipes(ctx, int64(1), 0).
					Return(nil)

				inventoryMock.EXPECT().UpdateSuperLikes(ctx, int64(1), 0).
					Return(nil)

				return &matcherServiceModule{
					cfg:              config,
					profileRepo:      profileMock,
					inventoryRepo:    inventoryMock,
					matchedRepo:      matchedMock,
					notificationRepo: notificationMock,
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := tt.mock()
			if err := m.SuperLike(tt.args.ctx, tt.args.userID, tt.args.targetUserID); (err != nil) != tt.wantErr {
				t.Errorf("matcherServiceModule.Like() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_matcherServiceModule_Dislike(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

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

	type args struct {
		ctx          context.Context
		userID       int64
		targetUserID int64
	}

	failedArgs := args{
		ctx:          ctx,
		userID:       1,
		targetUserID: 1,
	}

	successArgs := args{
		ctx:          ctx,
		userID:       1,
		targetUserID: 2,
	}

	tests := []struct {
		name    string
		mock    func() *matcherServiceModule
		args    args
		wantErr bool
	}{
		{
			name: "failed target user id is same with user id",
			mock: func() *matcherServiceModule {
				return &matcherServiceModule{
					cfg:              config,
					profileRepo:      nil,
					inventoryRepo:    nil,
					matchedRepo:      nil,
					notificationRepo: nil,
				}
			},
			args:    failedArgs,
			wantErr: true,
		},
		{
			name: "failed get inventory by user id",
			mock: func() *matcherServiceModule {
				inventoryMock := NewMockInventoryRepoInterface(ctrl)
				inventoryMock.EXPECT().GetByUserId(ctx, int64(1)).
					Return(nil, errors.New("test"))

				return &matcherServiceModule{
					cfg:              config,
					profileRepo:      nil,
					inventoryRepo:    inventoryMock,
					matchedRepo:      nil,
					notificationRepo: nil,
				}
			},
			args:    successArgs,
			wantErr: true,
		},
		{
			name: "failed insufficient swipes",
			mock: func() *matcherServiceModule {
				inventoryMock := NewMockInventoryRepoInterface(ctrl)
				inventoryMock.EXPECT().GetByUserId(ctx, int64(1)).
					Return(&domain.Inventory{
						Swipes: 0,
					}, nil)

				return &matcherServiceModule{
					cfg:              config,
					profileRepo:      nil,
					inventoryRepo:    inventoryMock,
					matchedRepo:      nil,
					notificationRepo: nil,
				}
			},
			args:    successArgs,
			wantErr: true,
		},
		{
			name: "failed update swipes",
			mock: func() *matcherServiceModule {
				inventoryMock := NewMockInventoryRepoInterface(ctrl)
				inventoryMock.EXPECT().GetByUserId(ctx, int64(1)).
					Return(&domain.Inventory{
						Swipes: 1,
					}, nil)

				inventoryMock.EXPECT().UpdateSwipes(ctx, int64(1), 0).
					Return(errors.New("test"))

				return &matcherServiceModule{
					cfg:              config,
					profileRepo:      nil,
					inventoryRepo:    inventoryMock,
					matchedRepo:      nil,
					notificationRepo: nil,
				}
			},
			args:    successArgs,
			wantErr: true,
		},
		{
			name: "success",
			mock: func() *matcherServiceModule {
				inventoryMock := NewMockInventoryRepoInterface(ctrl)
				inventoryMock.EXPECT().GetByUserId(ctx, int64(1)).
					Return(&domain.Inventory{
						Swipes: 1,
					}, nil)

				inventoryMock.EXPECT().UpdateSwipes(ctx, int64(1), 0).
					Return(nil)

				return &matcherServiceModule{
					cfg:              config,
					profileRepo:      nil,
					inventoryRepo:    inventoryMock,
					matchedRepo:      nil,
					notificationRepo: nil,
				}
			},
			args:    successArgs,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := tt.mock()
			if err := m.Dislike(tt.args.ctx, tt.args.userID, tt.args.targetUserID); (err != nil) != tt.wantErr {
				t.Errorf("matcherServiceModule.Dislike() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
