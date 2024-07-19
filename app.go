package main

import (
	"flag"
	"log"

	"github.com/gofiber/fiber/v2"
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

var (
	port = flag.String("port", ":3000", "Port to listen on")
	prod = flag.Bool("prod", false, "Enable prefork in Production")
)

func main() {
	// Parse command-line flags
	flag.Parse()

	// Create fiber app
	app := fiber.New(fiber.Config{
		// Prefork: *prod, // go run app.go -prod
	})

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

	// app.Use(cors.New(cors.Config{
	// 	AllowOrigins:     *coreSetting.AppConfig.Origin,
	// 	AllowCredentials: true,
	// 	AllowHeaders:     "Origin, Content-Type, Accept, X-Requested-With, X-HTTP-Method-Override, access-control-allow-origin, access-control-allow-headers",
	// }))

	// Listen on port 3000
	log.Fatal(app.Listen(*port)) // go run app.go -port=:3000
}
