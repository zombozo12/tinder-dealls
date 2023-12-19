package main

import (
	"fmt"
	"github.com/go-faker/faker/v4"
	"github.com/go-playground/validator/v10"
	"github.com/goccy/go-json"
	log "github.com/sirupsen/logrus"
	"github.com/zombozo12/tinder-dealls/domain"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"math/rand"
	"os"
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

	interestIn := []string{"male", "female"}
	gender := []string{"male", "female"}

	for i := 0; i < 100; i++ {
		user := domain.User{}
		if err := faker.FakeData(&user); err != nil {
			log.Panicf("Failed to fake user data: %s", err)
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
		if err != nil {
			log.Panicf("Failed to hash password: %s", err)
		}

		user.Password = string(hashedPassword)

		if result := db.Table("users").Create(&user); result.Error != nil {
			log.Panicf("Failed to insert user data: %s", err)
		}

		log.Printf("User: %v", user)

		profile := domain.Profile{}
		if err := faker.FakeData(&profile); err != nil {
			log.Panicf("Failed to fake profile data: %s", err)
		}

		profile.UserID = user.ID
		profile.InterestIn = randomPickerStr(interestIn)
		profile.Gender = randomPickerStr(gender)

		if result := db.Table("profile").Create(&profile); result.Error != nil {
			log.Panicf("Failed to insert profile data: %s", err)
		}

		log.Printf("Profile: %v", profile)

		inventory := domain.Inventory{
			UserID:     user.ID,
			Swipes:     rand.Int63n(11),
			Likes:      rand.Int63n(101),
			SuperLikes: rand.Int63n(101),
		}
		if result := db.Table("inventory").Create(&inventory); result.Error != nil {
			log.Panicf("Failed to insert inventory data: %s", err)
		}

		log.Printf("Inventory: %v", inventory)
	}
}

func randomPickerStr(data []string) string {
	return data[rand.Intn(len(data))]

}
