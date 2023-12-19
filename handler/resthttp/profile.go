package resthttp

import (
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
	"github.com/zombozo12/tinder-dealls/domain"
	"time"
)

type ProfileHandlerModule struct {
	cfg            *domain.Config
	profileService ProfileService
}

func NewProfileHandlerModule(cfg *domain.Config, profileService ProfileService) *ProfileHandlerModule {
	return &ProfileHandlerModule{
		cfg:            cfg,
		profileService: profileService,
	}
}

func (m ProfileHandlerModule) create(ctx *fiber.Ctx) error {
	startTime := time.Now()
	response := newResponse(ctx, startTime)
	tags := make(log.Fields)

	defer func() {
		tags["name"] = "handler.http.profile.create"
		tags["elapsed_time"] = time.Since(startTime).String()
		tags["request_id"] = ctx.Locals("requestid").(string)
		log.WithFields(tags).Debug()
	}()

	jwtUser, err := domain.ExtractUserClaims(ctx, m.cfg.JWT.Key)
	if err != nil {
		tags["error"] = "failed extracting user claims"
		tags["actual_error"] = err.Error()
		return response.setErrorResponse(fiber.StatusInternalServerError, "failed extracting user claims")
	}

	var req domain.ProfileRequest
	if err := ctx.BodyParser(&req); err != nil {
		tags["error"] = "failed parsing request"
		tags["actual_error"] = err.Error()
		return response.setErrorResponse(fiber.StatusBadRequest, "failed parsing request")
	}

	if err := m.profileService.Create(ctx.Context(), jwtUser.ID, req); err != nil {
		tags["error"] = "failed creating profile"
		tags["actual_error"] = err.Error()
		return response.setErrorResponse(fiber.StatusInternalServerError, "failed creating profile")
	}

	tags["status"] = "success"
	return response.setOKResponse(map[string]interface{}{"message": "profile created successfully"})
}

func (m ProfileHandlerModule) updateProfilePic(ctx *fiber.Ctx) error {
	startTime := time.Now()
	response := newResponse(ctx, startTime)
	tags := make(log.Fields)

	defer func() {
		tags["name"] = "handler.http.profile.update_profile_pic"
		tags["elapsed_time"] = time.Since(startTime).String()
		tags["request_id"] = ctx.Locals("requestid").(string)
		log.WithFields(tags).Debug()
	}()

	jwtUser, err := domain.ExtractUserClaims(ctx, m.cfg.JWT.Key)
	if err != nil {
		tags["error"] = "failed extracting user claims"
		tags["actual_error"] = err.Error()
		return response.setErrorResponse(fiber.StatusInternalServerError, "failed extracting user claims")
	}

	var req domain.UpdateProfilePicRequest
	if err := ctx.BodyParser(&req); err != nil {
		tags["error"] = "failed parsing request"
		tags["actual_error"] = err.Error()
		return response.setErrorResponse(fiber.StatusBadRequest, "failed parsing request")
	}

	if err := m.profileService.UpdateProfilePic(ctx.Context(), jwtUser.ID, req); err != nil {
		tags["error"] = "failed updating profile pic"
		tags["actual_error"] = err.Error()
		return response.setErrorResponse(fiber.StatusInternalServerError, "failed updating profile pic")
	}

	tags["status"] = "success"
	return response.setOKResponse(map[string]interface{}{"message": "profile pic updated successfully"})
}

func (m ProfileHandlerModule) updateProfile(ctx *fiber.Ctx) error {
	startTime := time.Now()
	response := newResponse(ctx, startTime)
	tags := make(log.Fields)

	defer func() {
		tags["name"] = "handler.http.profile.update_profile"
		tags["elapsed_time"] = time.Since(startTime).String()
		tags["request_id"] = ctx.Locals("requestid").(string)
		log.WithFields(tags).Debug()
	}()

	jwtUser, err := domain.ExtractUserClaims(ctx, m.cfg.JWT.Key)
	if err != nil {
		tags["error"] = "failed extracting user claims"
		tags["actual_error"] = err.Error()
		return response.setErrorResponse(fiber.StatusInternalServerError, "failed extracting user claims")
	}

	var req domain.ProfileRequest
	if err := ctx.BodyParser(&req); err != nil {
		tags["error"] = "failed parsing request"
		tags["actual_error"] = err.Error()
		return response.setErrorResponse(fiber.StatusBadRequest, "failed parsing request")
	}

	if err := m.profileService.UpdateProfile(ctx.Context(), jwtUser.ID, req); err != nil {
		tags["error"] = "failed updating profile"
		tags["actual_error"] = err.Error()
		return response.setErrorResponse(fiber.StatusInternalServerError, "failed updating profile")
	}

	tags["status"] = "success"
	return response.setOKResponse(map[string]interface{}{"message": "profile updated successfully"})
}
