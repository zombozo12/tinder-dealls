package resthttp

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/zombozo12/tinder-dealls/domain"
)

type RouteDependencies struct {
	Cfg            *domain.Config
	Auth           AuthService
	Profile        ProfileService
	Recommendation RecommendationService
	Matcher        MatcherService
}

func NewRouter(app *fiber.App, dep RouteDependencies) {
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(requestid.New())

	authMiddleware := authenticate(dep.Cfg.Key)

	authHandler := NewAuthHandlerModule(dep.Auth)
	profileHandler := NewProfileHandlerModule(dep.Cfg, dep.Profile)
	recommendationHandler := NewRecommendationHandlerModule(dep.Cfg, dep.Recommendation)
	matcherHandler := NewMatcherHandlerModule(dep.Cfg, dep.Matcher)

	// Set global prefix to /api
	api := app.Group("/api")

	// Set prefix to /api/auth
	auth := api.Group("/auth")
	auth.Post("/in", authHandler.login)
	auth.Post("/up", authHandler.register)

	// Set prefix to /api/profile
	profile := api.Group("/profile").Use(authMiddleware)
	profile.Post("/create", profileHandler.create)
	profile.Put("/pic", profileHandler.updateProfilePic)
	profile.Put("/update", profileHandler.updateProfile)

	// Set prefix to /api/recommendation
	recommendation := api.Group("/recommendation").Use(authMiddleware)
	recommendation.Get("/get", recommendationHandler.getRecommendation)

	// Set prefix to /api/matcher
	matcher := api.Group("/matcher").Use(authMiddleware)
	matcher.Post("/like", matcherHandler.like)
	matcher.Post("/superlike", matcherHandler.superLike)
	matcher.Post("/dislike", matcherHandler.dislike)
}
