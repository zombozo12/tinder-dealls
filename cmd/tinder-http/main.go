package main

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
	"github.com/zombozo12/tinder-dealls/domain"
	"github.com/zombozo12/tinder-dealls/handler/resthttp"
	"github.com/zombozo12/tinder-dealls/repository/auth"
	"github.com/zombozo12/tinder-dealls/repository/inventory"
	"github.com/zombozo12/tinder-dealls/repository/matched"
	"github.com/zombozo12/tinder-dealls/repository/notification"
	"github.com/zombozo12/tinder-dealls/repository/profile"
	"github.com/zombozo12/tinder-dealls/repository/rds"
	"github.com/zombozo12/tinder-dealls/services"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

func main() {
	log.SetFormatter(&log.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})

	log.SetLevel(log.DebugLevel)

	currentDir, err := os.Getwd()
	if err != nil {
		log.Panicf("Failed to get current directory path: %s", err)
	}

	configPath := fmt.Sprintf("%s/config.json", currentDir)
	configFile, err := os.ReadFile(configPath)
	if err != nil {
		log.Panicf("Failed to read config file: %s", err)
	}

	var config *domain.Config
	err = json.Unmarshal(configFile, &config)
	if err != nil {
		log.Panicf("Failed to unmarshal config file: %s", err)
	}

	validate := validator.New()
	err = validate.Struct(config)
	if err != nil {
		log.Panicf("Failed to validate config file: %s", err)
	}

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", config.Database.Host, config.Database.Port, config.Database.Username, config.Database.Password, config.Database.Name)
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})
	if err != nil {
		log.Panicf("Failed to connect to database: %s", err)
	}

	app := fiber.New(fiber.Config{
		Prefork:       false,
		ServerHeader:  "tinder",
		StrictRouting: true,
		CaseSensitive: true,
		AppName:       "Tinder Dealls",
		JSONEncoder:   json.Marshal,
		JSONDecoder:   json.Unmarshal,
	})

	redisClient := redis.NewClient(&redis.Options{
		Addr:     config.Redis.Host + ":" + strconv.Itoa(config.Redis.Port),
		Password: config.Redis.Password,
		DB:       0,
	})

	authRepo := auth.New(db, config)
	profileRepo := profile.New(db, config)
	inventoryRepo := inventory.New(db, config)
	notificationRepo := notification.New(db, config)
	redisRepo := rds.New(redisClient, config)
	matchedRepo := matched.New(db, config)

	// Setting up services
	authService, err := services.NewAuthService(config, db, authRepo, inventoryRepo, profileRepo)
	if err != nil {
		log.Panicf("Failed to setup auth service: %s", err)
	}

	profileService, err := services.NewProfileService(config, db, profileRepo)
	if err != nil {
		log.Panicf("Failed to setup profile service: %s", err)
	}

	matcherService, err := services.NewMatcherService(config,
		profileRepo, inventoryRepo, matchedRepo, notificationRepo)
	if err != nil {
		log.Panicf("Failed to setup matcher service: %s", err)
	}

	recommendationService, err := services.NewRecommendationService(config, db, redisRepo, profileRepo)
	if err != nil {
		log.Panicf("Failed to setup recommendation service: %s", err)
	}

	// Setting up router
	resthttp.NewRouter(app, resthttp.RouteDependencies{
		Cfg:            config,
		Auth:           authService,
		Profile:        profileService,
		Matcher:        matcherService,
		Recommendation: recommendationService,
	})

	// Setting up graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		_ = <-c
		log.Info("Gracefully shutting down...")
		_ = app.Shutdown()
	}()

	if err := app.Listen(fmt.Sprintf(":%d", config.Server.Port)); err != nil {
		log.Panicf("Failed to start server: %s", err)
	}

	log.Info("Running cleanup tasks...")
}
