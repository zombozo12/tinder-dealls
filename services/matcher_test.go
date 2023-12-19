package services

import (
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
