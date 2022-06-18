package auth

import (
	"context"
	"embed"
	"io/fs"
	"net/http"

	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
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
	config.InitCoreConfig(&coreSetting.AppConfig)
	coreSetting.AppConfig.PayloadSecret = &secrets.PayloadSecret
	coreSetting.AppConfig.PublicKey = &secrets.KeyPub
	coreSetting.AppConfig.PrivateKey = &secrets.Key
	coreSetting.AppConfig.RefEmailPass = &secrets.RefEmailPass
	coreSetting.AppConfig.RecaptchaKey = &secrets.RecaptchaKey
	coreSetting.AppConfig.MongoDBHost = &secrets.MongoHost
	coreSetting.AppConfig.Database = &secrets.MongoDatabase

	// Init auth micro config
	config.InitAuthConfig(&authSetting.AuthConfig)
	authSetting.AuthConfig.OAuthClientSecret = secrets.TSClientSecret
	authSetting.AuthConfig.AdminUsername = secrets.AdminUsername
	authSetting.AuthConfig.AdminPassword = secrets.AdminPassword

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
	app.Use(cors.New(cors.Config{
		AllowOrigins:     *coreSetting.AppConfig.Origin,
		AllowCredentials: true,
		AllowHeaders:     "Origin, Content-Type, Accept, X-Requested-With, X-HTTP-Method-Override, access-control-allow-origin, access-control-allow-headers",
	}))
	// app.Use(func (c *fiber.Ctx) error {
	// 	c.Locals("app", "auth")
	// 	return nil
	// })
	router.SetupRoutes(app)
}

// Auth handler
//
//encore:api public raw path=/auth/*p1
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
