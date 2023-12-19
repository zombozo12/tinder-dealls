package resthttp

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"regexp"
	"strings"
	"time"
)

type response struct {
	Context   *fiber.Ctx  `json:"-"`
	IsError   bool        `json:"is_error"`
	Data      interface{} `json:"data,omitempty"`
	Elapsed   string      `json:"elapsed_time"`
	RequestID string      `json:"request_id"`
	Start     time.Time   `json:"-"`
}

func newResponse(ctx *fiber.Ctx, start time.Time) response {
	return response{
		Context:   ctx,
		Start:     start,
		RequestID: ctx.Locals("requestid").(string),
	}
}

func (r *response) setOKResponse(data interface{}) error {
	r.Data = data
	r.Elapsed = time.Since(r.Start).String()
	r.IsError = false

	return r.Context.Status(fiber.StatusOK).JSON(r)
}

func (r *response) setErrorResponse(statusCode int, err string) error {
	r.Data = map[string]string{
		"error": err,
	}
	r.Elapsed = time.Since(r.Start).String()
	r.IsError = true

	return r.Context.Status(statusCode).JSON(r)
}

func (r *response) setErrorValidationResponse(err error) error {
	r.Elapsed = time.Since(r.Start).String()
	r.IsError = true

	var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

	errors := make(map[string]string)
	for _, v := range err.(validator.ValidationErrors) {
		var sb strings.Builder

		snake := matchFirstCap.ReplaceAllString(v.Field(), "${1}_${2}")
		snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
		field := strings.ToLower(snake)

		sb.WriteString(field + " should be " + v.Tag() + " ")

		if v.Param() != "" {
			sb.WriteString(v.Param())
		}

		errors[field] = strings.TrimSpace(sb.String())
	}

	r.Data = map[string]map[string]string{
		"error": errors,
	}

	return r.Context.Status(fiber.StatusBadRequest).JSON(r)
}
