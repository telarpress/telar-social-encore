package auth

import (
	"context"
	"embed"
	"io/fs"
	"net/http"

	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/template/html"
	coreSetting "github.com/red-gold/telar-core/config"
	"github.com/red-gold/telar-core/pkg/log"
	authSetting "github.com/red-gold/telar-web/micros/auth/config"
	"github.com/red-gold/telar-web/micros/auth/database"
	"github.com/red-gold/telar-web/micros/auth/router"
	"github.com/telarpress/telar-social-encore/config"
)

// Cache state
var app *fiber.App

//go:embed all:views
var viewsFS embed.FS

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
	config.InitCoreConfig("auth", &coreSetting.AppConfig, allSecrets)
	println(*coreSetting.AppConfig.MongoDBHost)
	// Init auth micro config
	config.InitAuthConfig(&authSetting.AuthConfig, allSecrets)
	println(authSetting.AuthConfig.AdminPassword)

	// Initialize app
	dirFS, err := fs.Sub(viewsFS, "views")
	if err != nil {
		log.Error(err.Error())
	}

	engine := html.NewFileSystem(http.FS(dirFS), ".html")
	app = fiber.New(fiber.Config{
		Views: engine,
	})

	app.Use(recover.New())
	app.Use(requestid.New())
	app.Use(logger.New(
		logger.Config{
			Format: "[${time}] ${status} - ${latency} ${method} ${path} - ${header:}\n​",
		},
	))

	router.SetupRoutes(app)
}

// Mount the app to the parent app
func Mount(route string, parentApp *fiber.App) {
	app.Use(func(c *fiber.Ctx) error {
		return ConnectDatabase(c.Context())
	})
	parentApp.Mount(route, app)

}

// Connect to database
func ConnectDatabase(ctx context.Context) error {

	// Connect
	if database.Db == nil {
		startErr := database.Connect(ctx)
		if startErr != nil {
			log.Error("Error connect to database: %s", startErr.Error())
			return startErr
		}
	}
	return nil
}

// Auth handler
//
//encore:api public raw path=/auth/*p1
func Handle(w http.ResponseWriter, r *http.Request) {
	// Remove base url from request path
	RemoveBaseURLFromRequest(r)

	ctx := context.Background()

	// Connect
	startErr := ConnectDatabase(ctx)
	if startErr != nil {
		log.Error("Error startup: %s", startErr.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(startErr.Error()))
	}

	adaptor.FiberApp(app)(w, r)
}
