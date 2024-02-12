package test

import (
	"github.com/casbin/casbin/v2"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/n9mi/go-course-app/internal/config"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var (
	viperConfig *viper.Viper
	log         *logrus.Logger
	app         *fiber.App
	db          *gorm.DB
	redisClient *redis.Client
	validate    *validator.Validate
	enforcer    *casbin.Enforcer
)

type UserValidData struct {
	ID    string
	Name  string
	Email string
	Roles []string
	Token string
}

var validAdminData = UserValidData{Name: "Admin 1", Email: "admin1@mail.com"}
var validUserData = UserValidData{Name: "User 1", Email: "user1@mail.com"}

type TestSchema map[string]interface{}

type TestResponse[T any] struct {
	Data     T        `json:"data"`
	Messages []string `json:"messages"`
}

func init() {
	viperConfig = config.NewViper()
	log = config.NewLogger(viperConfig)
	app = config.NewFiber(viperConfig)
	db = config.NewDatabase(viperConfig, log)
	redisClient = config.NewRedisClient(viperConfig)
	validate = config.NewValidator(viperConfig)
	enforcer = newTestEnforcer(log)

	configBootstrap := &config.ConfigBootstrap{
		ViperConfig: viperConfig,
		App:         app,
		DB:          db,
		Validate:    validate,
		RedisClient: redisClient,
		Log:         log,
		Enforcer:    enforcer,
	}
	config.Bootstrap(configBootstrap)
}
