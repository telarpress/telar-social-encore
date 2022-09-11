package votes

import (
	"context"
	"net/http"

	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	coreSetting "github.com/red-gold/telar-core/config"
	"github.com/red-gold/telar-core/pkg/log"
	"github.com/red-gold/ts-serverless/micros/votes/database"
	"github.com/red-gold/ts-serverless/micros/votes/router"
	"github.com/telarpress/telar-social-encore/config"
)

// Cache state
var app *fiber.App

// Secrets
var secrets struct {
	AdminUsername  string
	AdminPassword  string
	MongoHost      string
	MongoDatabase  string
	PhoneAuthId    string
	PhoneAuthToken string
	Key            string
	KeyPub         string
	RefEmailPass   string
	PayloadSecret  string
	ServiceAccount string
	TSClientSecret string
	RecaptchaKey   string
}

// init
func init() {

	// Init core config
	var allSecrets = &config.AllSecrets{
		AdminUsername:  secrets.AdminPassword,
		AdminPassword:  secrets.AdminPassword,
		MongoHost:      secrets.MongoHost,
		MongoDatabase:  secrets.MongoDatabase,
		PhoneAuthId:    secrets.PhoneAuthId,
		PhoneAuthToken: secrets.PhoneAuthToken,
		Key:            secrets.Key,
		KeyPub:         secrets.KeyPub,
		RefEmailPass:   secrets.RefEmailPass,
		PayloadSecret:  secrets.PayloadSecret,
		ServiceAccount: secrets.ServiceAccount,
		TSClientSecret: secrets.TSClientSecret,
		RecaptchaKey:   secrets.RecaptchaKey,
	}
	config.InitCoreConfig("votes", &coreSetting.AppConfig, allSecrets)

	// Initialize app
	app = fiber.New()

	app.Use(recover.New())
	app.Use(requestid.New())
	app.Use(logger.New(
		logger.Config{
			Format: "[${time}] ${status} - ${latency} ${method} ${path} - ${header:}\n​",
		},
	))
	// app.Use(cors.New(cors.Config{
	// 	AllowOrigins:     *coreSetting.AppConfig.Origin,
	// 	AllowCredentials: true,
	// 	AllowHeaders:     "Origin, Content-Type, Accept, X-Requested-With, X-HTTP-Method-Override, access-control-allow-origin, access-control-allow-headers",
	// }))
	router.SetupRoutes(app)
}

// Votes handler
//
//encore:api public raw path=/votes/*p1
func Handle(w http.ResponseWriter, r *http.Request) {
	// Remove base url from request path
	RemoveBaseURLFromRequest(r)

	ctx := context.Background()

	// Connect
	if database.Db == nil {
		startErr := database.Connect(ctx)
		if startErr != nil {
			log.Error("Error startup: %s", startErr.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(startErr.Error()))
		}
	}

	adaptor.FiberApp(app)(w, r)
}
