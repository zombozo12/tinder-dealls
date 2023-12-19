package services

import (
	"context"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/zombozo12/tinder-dealls/domain"
	"time"
)

type matcherServiceModule struct {
	cfg              *domain.Config
	profileRepo      ProfileRepoInterface
	inventoryRepo    InventoryRepoInterface
	matchedRepo      MatchedRepoInterface
	notificationRepo NotificationRepoInterface
}

type MatcherServiceInterface interface {
	Like(ctx context.Context, userID int64, targetUserID int64) error
	SuperLike(ctx context.Context, userID int64, targetUserID int64) error
	Dislike(ctx context.Context, userID int64, targetUserID int64) error
}

func NewMatcherService(cfg *domain.Config,
	profileRepo ProfileRepoInterface,
	inventoryRepo InventoryRepoInterface,
	matchedRepo MatchedRepoInterface,
	notificationRepo NotificationRepoInterface) (MatcherServiceInterface, error) {
	return &matcherServiceModule{
		cfg:              cfg,
		profileRepo:      profileRepo,
		inventoryRepo:    inventoryRepo,
		matchedRepo:      matchedRepo,
		notificationRepo: notificationRepo,
	}, nil
}

func (m matcherServiceModule) Like(ctx context.Context, userID int64, targetUserID int64) error {
	startTime := time.Now()
	tags := make(log.Fields)

	defer func() {
		tags["name"] = "service.matcher.like"
		tags["elapsed_time"] = time.Since(startTime).String()
		tags["request_id"] = ctx.Value("requestid").(string)
		log.WithFields(tags).Debug()
	}()

	if targetUserID == userID {
		tags["error"] = "cannot like yourself"
		tags["status"] = "error"
		return errors.New("cannot like yourself")
	}

	inventory, err := m.inventoryRepo.GetByUserId(ctx, userID)
	if err != nil {
		tags["error"] = "failed get inventory"
		tags["status"] = "error"
		return err
	}

	if inventory.Swipes < 1 {
		tags["error"] = "insufficient swipes"
		tags["status"] = "error"
		return errors.New("insufficient swipes")
	}

	if inventory.Likes < 1 {
		tags["error"] = "insufficient likes"
		tags["status"] = "error"
		return errors.New("insufficient likes")
	}

	targetProfile, err := m.profileRepo.GetProfile(ctx, userID)
	if err != nil {
		tags["error"] = "failed get profile"
		tags["status"] = "error"
		return err
	}

	if targetProfile == nil {
		tags["error"] = "profile not found"
		tags["status"] = "error"
		return errors.New("profile not found")
	}

	isExists, err := m.matchedRepo.IsExists(ctx, domain.MatchRequest{
		UserID:       userID,
		TargetUserID: targetUserID,
	})
	if err != nil {
		tags["error"] = "failed check matched"
		tags["status"] = "error"
		return err
	}

	if isExists {
		tags["error"] = "already matched"
		tags["status"] = "warning"
		return errors.New("already matched")
	}

	if err := m.matchedRepo.Create(ctx, domain.MatchRequest{
		UserID:       userID,
		TargetUserID: targetUserID,
	}); err != nil {
		tags["error"] = "failed create matched"
		tags["status"] = "error"
		return err
	}

	isMatched, err := m.matchedRepo.IsMatched(ctx, domain.MatchRequest{
		UserID:       targetUserID,
		TargetUserID: userID,
	})
	if err != nil {
		tags["error"] = "failed check matched"
		tags["status"] = "error"
		return err
	}

	if isMatched {
		if err := m.notificationRepo.Create(ctx, domain.NotificationRequest{
			UserID:  targetUserID,
			Message: fmt.Sprintf("You have a new match with %s", targetProfile.Name),
		}); err != nil {
			tags["error"] = "failed create notification"
			tags["status"] = "error"
			return err
		}
		return nil
	}

	if err := m.inventoryRepo.UpdateSwipes(ctx, userID, int(inventory.Swipes)-1); err != nil {
		tags["error"] = "failed update swipes"
		tags["status"] = "error"
		return err
	}

	if err := m.inventoryRepo.UpdateLikes(ctx, userID, int(inventory.Likes)-1); err != nil {
		tags["error"] = "failed update likes"
		tags["status"] = "error"
		return err
	}

	tags["status"] = "success"
	return nil
}

func (m matcherServiceModule) SuperLike(ctx context.Context, userID int64, targetUserID int64) error {
	startTime := time.Now()
	tags := make(log.Fields)

	defer func() {
		tags["name"] = "service.matcher.super_like"
		tags["elapsed_time"] = time.Since(startTime).String()
		tags["request_id"] = ctx.Value("requestid").(string)
		log.WithFields(tags).Debug()
	}()

	if targetUserID == userID {
		tags["error"] = "cannot super like yourself"
		tags["status"] = "error"
		return errors.New("cannot super like yourself")
	}

	inventory, err := m.inventoryRepo.GetByUserId(ctx, userID)
	if err != nil {
		tags["error"] = "failed get inventory"
		tags["status"] = "error"
		return err
	}

	if inventory.Swipes < 1 {
		tags["error"] = "insufficient swipes"
		tags["status"] = "error"
		return errors.New("insufficient swipes")
	}

	if inventory.SuperLikes < 1 {
		tags["error"] = "insufficient super likes"
		tags["status"] = "error"
		return errors.New("insufficient super likes")
	}

	targetProfile, err := m.profileRepo.GetProfile(ctx, userID)
	if err != nil {
		tags["error"] = "failed get profile"
		tags["status"] = "error"
		return err
	}

	if targetProfile == nil {
		tags["error"] = "profile not found"
		tags["status"] = "error"
		return errors.New("profile not found")
	}

	isExists, err := m.matchedRepo.IsExists(ctx, domain.MatchRequest{
		UserID:       userID,
		TargetUserID: targetUserID,
	})
	if err != nil {
		tags["error"] = "failed check exists"
		tags["status"] = "error"
		return err
	}

	if isExists {
		tags["error"] = "already matched"
		tags["status"] = "warning"
		return errors.New("already matched")
	}

	if err := m.matchedRepo.Create(ctx, domain.MatchRequest{
		UserID:       userID,
		TargetUserID: targetUserID,
	}); err != nil {
		tags["error"] = "failed create matched"
		tags["status"] = "error"
		return err
	}

	isMatched, err := m.matchedRepo.IsMatched(ctx, domain.MatchRequest{
		UserID:       targetUserID,
		TargetUserID: userID,
	})
	if err != nil {
		tags["error"] = "failed check matched"
		tags["status"] = "error"
		return err
	}

	if isMatched {
		if err := m.notificationRepo.Create(ctx, domain.NotificationRequest{
			UserID:  targetUserID,
			Message: fmt.Sprintf("You have a new match with %s", targetProfile.Name),
		}); err != nil {
			tags["error"] = "failed create notification"
			tags["status"] = "error"
			return err
		}
		return nil
	}

	if err := m.notificationRepo.Create(ctx, domain.NotificationRequest{
		UserID:  targetUserID,
		Message: fmt.Sprintf("%s super liked you", targetProfile.Name),
	}); err != nil {
		tags["error"] = "failed create notification"
		tags["status"] = "error"
		return err
	}

	if err := m.inventoryRepo.UpdateSwipes(ctx, userID, int(inventory.Swipes)-1); err != nil {
		tags["error"] = "failed update swipes"
		tags["status"] = "error"
		return err
	}

	if err := m.inventoryRepo.UpdateSuperLikes(ctx, userID, int(inventory.SuperLikes)-1); err != nil {
		tags["error"] = "failed update super likes"
		tags["status"] = "error"
		return err
	}

	tags["status"] = "success"
	return nil
}

func (m matcherServiceModule) Dislike(ctx context.Context, userID int64, targetUserID int64) error {
	startTime := time.Now()
	tags := make(log.Fields)

	defer func() {
		tags["name"] = "service.matcher.dislike"
		tags["elapsed_time"] = time.Since(startTime).String()
		tags["request_id"] = ctx.Value("requestid").(string)
		log.WithFields(tags).Debug()
	}()

	if targetUserID == userID {
		tags["error"] = "cannot dislike yourself"
		tags["status"] = "error"
		return errors.New("cannot dislike yourself")
	}

	inventory, err := m.inventoryRepo.GetByUserId(ctx, userID)
	if err != nil {
		tags["error"] = "failed get inventory"
		tags["status"] = "error"
		return err
	}

	if inventory.Swipes < 1 {
		tags["error"] = "insufficient swipes"
		tags["status"] = "error"
		return errors.New("insufficient swipes")
	}

	if err := m.inventoryRepo.UpdateSwipes(ctx, userID, int(inventory.Swipes)-1); err != nil {
		tags["error"] = "failed update swipes"
		tags["status"] = "error"
		return err
	}

	tags["status"] = "success"
	return nil
}
