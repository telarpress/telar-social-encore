package main

import (
	"embed"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/html"
	coreSetting "github.com/red-gold/telar-core/config"
	actionsMicro "github.com/telarpress/telar-social-encore/actions"
	adminMicro "github.com/telarpress/telar-social-encore/admin"
	authMicro "github.com/telarpress/telar-social-encore/auth"
	circlesMicro "github.com/telarpress/telar-social-encore/circles"
	commentsMicro "github.com/telarpress/telar-social-encore/comments"
	mediaMicro "github.com/telarpress/telar-social-encore/gallery"
	notificationsMicro "github.com/telarpress/telar-social-encore/notifications"
	postsMicro "github.com/telarpress/telar-social-encore/posts"
	profileMicro "github.com/telarpress/telar-social-encore/profile"
	settingMicro "github.com/telarpress/telar-social-encore/setting"
	storageMicro "github.com/telarpress/telar-social-encore/storage"
	userRelsMicro "github.com/telarpress/telar-social-encore/user-rels"
	vangMicro "github.com/telarpress/telar-social-encore/vang"
	votesMicro "github.com/telarpress/telar-social-encore/votes"
)

//go:embed all:auth/views
var viewsFS embed.FS

var (
	port = flag.String("port", ":3000", "Port to listen on")
	prod = flag.Bool("prod", false, "Enable prefork in Production")
)

func main() {
	// Parse command-line flags
	flag.Parse()
	dirFS, err := fs.Sub(viewsFS, "views")
	if err != nil {
		log.Fatal(err.Error())
	}

	engine := html.NewFileSystem(http.FS(dirFS), ".html")

	// Create fiber app
	app := fiber.New(fiber.Config{
		// Prefork: *prod, // go run app.go -prod
		Views: engine,
		// Override default error handler
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			// Status code defaults to 500
			code := fiber.StatusInternalServerError
			log.Print(err.Error())
			// Retrieve the custom status code if it's a *fiber.Error
			var e *fiber.Error
			if errors.As(err, &e) {
				code = e.Code
			}

			// Send custom error page
			err = ctx.Status(code).SendFile(fmt.Sprintf("./%d.html", code))
			if err != nil {
				// In case the SendFile fails
				return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
			}

			// Return from handler
			return nil
		},
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins:     *coreSetting.AppConfig.Origin,
		AllowCredentials: true,
		AllowHeaders:     "Origin, Authorization, Content-Type, Accept, X-Requested-With, X-HTTP-Method-Override, access-control-allow-credentials",
		AllowMethods:     "GET, POST, HEAD, PUT, DELETE, PATCH, OPTIONS",
	}))
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		// For more options, see the Config section
		Format: "${pid} ${locals:requestid} ${status} - ${method} ${path}​\n",
	}))

	actionsMicro.Mount("/actions", app)
	adminMicro.Mount("/admin", app)
	authMicro.Mount("/auth", app)
	circlesMicro.Mount("/circles", app)
	commentsMicro.Mount("/comments", app)
	mediaMicro.Mount("/media", app)
	notificationsMicro.Mount("/notifications", app)
	postsMicro.Mount("/posts", app)
	profileMicro.Mount("/profile", app)
	settingMicro.Mount("/setting", app)
	storageMicro.Mount("/storage", app)
	userRelsMicro.Mount("/user-rels", app)
	votesMicro.Mount("/votes", app)
	vangMicro.Mount("/vang", app)

	// Listen on port 3000
	log.Fatal(app.Listen(*port)) // go run app.go -port=:3000
}
