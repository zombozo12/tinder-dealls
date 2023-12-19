package resthttp

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
	"github.com/zombozo12/tinder-dealls/domain"
	"time"
)

type AuthHandlerModule struct {
	authService AuthService
}

func NewAuthHandlerModule(authService AuthService) *AuthHandlerModule {
	return &AuthHandlerModule{
		authService: authService,
	}
}

func (m AuthHandlerModule) login(ctx *fiber.Ctx) error {
	startTime := time.Now()
	response := newResponse(ctx, startTime)
	tags := make(log.Fields)

	defer func() {
		tags["name"] = "handler.http.auth.login"
		tags["elapsed_time"] = time.Since(startTime).String()
		tags["request_id"] = ctx.Locals("requestid").(string)
		log.WithFields(tags).Debug()
	}()

	var req domain.AuthRequest
	if err := ctx.BodyParser(&req); err != nil {
		tags["error"] = "failed parsing request"
		return response.setErrorResponse(fiber.StatusUnprocessableEntity, "failed parsing request")
	}

	validate := validator.New()
	err := validate.Struct(req)
	if err != nil {
		tags["error"] = "failed validating request"
		return response.setErrorValidationResponse(err)
	}

	res, err := m.authService.Login(ctx.Context(), req)
	if err != nil {
		tags["error"] = "failed login"
		return response.setErrorResponse(fiber.StatusUnauthorized, "failed login")
	}

	tags["status"] = "success"
	return response.setOKResponse(res)
}

func (m AuthHandlerModule) register(ctx *fiber.Ctx) error {
	startTime := time.Now()
	response := newResponse(ctx, startTime)
	tags := make(log.Fields)
	defer func() {
		tags["name"] = "handler.http.auth.register"
		tags["elapsed_time"] = time.Since(startTime).String()
		tags["request_id"] = ctx.Locals("requestid").(string)
		log.WithFields(tags).Debug()
	}()

	var req domain.AuthRequest
	if err := ctx.BodyParser(&req); err != nil {
		tags["error"] = "failed parsing request"
		return response.setErrorResponse(fiber.StatusUnprocessableEntity, "failed parsing request")
	}

	validate := validator.New()
	err := validate.Struct(req)
	if err != nil {
		tags["error"] = "failed validating request"
		return response.setErrorValidationResponse(err)
	}

	err = m.authService.Register(ctx.Context(), req)
	if err != nil {
		tags["error"] = "failed register"
		return response.setErrorResponse(fiber.StatusUnauthorized, "failed register")
	}

	tags["status"] = "success"
	return response.setOKResponse(map[string]string{
		"message": "success register",
	})
}
