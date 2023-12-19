package resthttp

import (
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
	"github.com/zombozo12/tinder-dealls/domain"
	"time"
)

type MatcherHandlerModule struct {
	cfg            *domain.Config
	matcherService MatcherService
}

func NewMatcherHandlerModule(cfg *domain.Config, matcherService MatcherService) *MatcherHandlerModule {
	return &MatcherHandlerModule{
		cfg:            cfg,
		matcherService: matcherService,
	}
}

func (h *MatcherHandlerModule) like(ctx *fiber.Ctx) error {
	startTime := time.Now()
	response := newResponse(ctx, startTime)
	tags := make(log.Fields)

	defer func() {
		tags["name"] = "handler.http.matcher.like"
		tags["elapsed_time"] = time.Since(startTime).String()
		tags["request_id"] = ctx.Locals("requestid").(string)
		log.WithFields(tags).Debug()
	}()

	jwtUser, err := domain.ExtractUserClaims(ctx, h.cfg.JWT.Key)
	if err != nil {
		tags["error"] = "failed extracting user claims"
		tags["actual_error"] = err.Error()
		return response.setErrorResponse(fiber.StatusInternalServerError, "failed extracting user claims")
	}

	var req domain.MatchRequest
	if err := ctx.BodyParser(&req); err != nil {
		tags["error"] = "failed parsing request"
		tags["actual_error"] = err.Error()
		return response.setErrorResponse(fiber.StatusBadRequest, "failed parsing request")
	}

	if err := h.matcherService.Like(ctx.Context(), jwtUser.ID, req.TargetUserID); err != nil {
		tags["error"] = "failed liking"
		tags["actual_error"] = err.Error()
		return response.setErrorResponse(fiber.StatusInternalServerError, "failed liking")
	}

	tags["status"] = "success"
	return response.setOKResponse(map[string]interface{}{"message": "liked successfully"})
}

func (h *MatcherHandlerModule) superLike(ctx *fiber.Ctx) error {
	startTime := time.Now()
	response := newResponse(ctx, startTime)
	tags := make(log.Fields)

	defer func() {
		tags["name"] = "handler.http.matcher.superLike"
		tags["elapsed_time"] = time.Since(startTime).String()
		tags["request_id"] = ctx.Locals("requestid").(string)
		log.WithFields(tags).Debug()
	}()

	jwtUser, err := domain.ExtractUserClaims(ctx, h.cfg.JWT.Key)
	if err != nil {
		tags["error"] = "failed extracting user claims"
		tags["actual_error"] = err.Error()
		return response.setErrorResponse(fiber.StatusInternalServerError, "failed extracting user claims")
	}

	var req domain.MatchRequest
	if err := ctx.BodyParser(&req); err != nil {
		tags["error"] = "failed parsing request"
		tags["actual_error"] = err.Error()
		return response.setErrorResponse(fiber.StatusBadRequest, "failed parsing request")
	}

	if err := h.matcherService.SuperLike(ctx.Context(), jwtUser.ID, req.TargetUserID); err != nil {
		tags["error"] = "failed super liking"
		tags["actual_error"] = err.Error()
		return response.setErrorResponse(fiber.StatusInternalServerError, "failed super liking")
	}

	tags["status"] = "success"
	return response.setOKResponse(map[string]interface{}{"message": "super liked successfully"})
}

func (h *MatcherHandlerModule) dislike(ctx *fiber.Ctx) error {
	startTime := time.Now()
	response := newResponse(ctx, startTime)
	tags := make(log.Fields)

	defer func() {
		tags["name"] = "handler.http.matcher.dislike"
		tags["elapsed_time"] = time.Since(startTime).String()
		tags["request_id"] = ctx.Locals("requestid").(string)
		log.WithFields(tags).Debug()
	}()

	jwtUser, err := domain.ExtractUserClaims(ctx, h.cfg.JWT.Key)
	if err != nil {
		tags["error"] = "failed extracting user claims"
		tags["actual_error"] = err.Error()
		return response.setErrorResponse(fiber.StatusInternalServerError, "failed extracting user claims")
	}

	var req domain.MatchRequest
	if err := ctx.BodyParser(&req); err != nil {
		tags["error"] = "failed parsing request"
		tags["actual_error"] = err.Error()
		return response.setErrorResponse(fiber.StatusBadRequest, "failed parsing request")
	}

	if err := h.matcherService.Dislike(ctx.Context(), jwtUser.ID, req.TargetUserID); err != nil {
		tags["error"] = "failed disliking"
		tags["actual_error"] = err.Error()
		return response.setErrorResponse(fiber.StatusInternalServerError, "failed disliking")
	}

	tags["status"] = "success"
	return response.setOKResponse(map[string]interface{}{"message": "disliked successfully"})
}
