package resthttp

import (
	"github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"time"
)

func authenticate(key string) fiber.Handler {
	return jwtware.New(jwtware.Config{
		SuccessHandler: func(ctx *fiber.Ctx) error {
			return ctx.Next()
		},
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			response := newResponse(ctx, time.Now())
			return response.setErrorResponse(fiber.StatusUnauthorized, "unauthorized")
		},
		SigningKey: jwtware.SigningKey{
			Key: []byte(key),
		},
	})
}
