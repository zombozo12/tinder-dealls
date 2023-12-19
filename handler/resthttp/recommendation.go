package resthttp

import (
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
	"github.com/zombozo12/tinder-dealls/domain"
	"time"
)

type RecommendationHandlerModule struct {
	cfg                   *domain.Config
	recommendationService RecommendationService
}

func NewRecommendationHandlerModule(cfg *domain.Config, recommendationService RecommendationService) *RecommendationHandlerModule {
	return &RecommendationHandlerModule{
		cfg:                   cfg,
		recommendationService: recommendationService,
	}
}

func (h *RecommendationHandlerModule) getRecommendation(ctx *fiber.Ctx) error {
	startTime := time.Now()
	response := newResponse(ctx, startTime)
	tags := make(log.Fields)
	defer func() {
		tags["name"] = "handler.http.recommendation.getRecommendation"
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

	recommendation, err := h.recommendationService.GetRecommendation(ctx.Context(), jwtUser.ID)
	if err != nil {
		tags["error"] = "failed getting recommendation"
		return response.setErrorResponse(fiber.StatusInternalServerError, "failed getting recommendation")
	}

	tags["status"] = "success"
	return response.setOKResponse(recommendation)
}
