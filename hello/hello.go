package hello

import (
	"net/http"

	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
)

// Cache state
var app *fiber.App


var secrets struct {
    GitHubAPIKey string 
}


// init
func init() {


	// Initialize app
	app = fiber.New()

	app.Get("/:name", func(c *fiber.Ctx) error {
		name := c.Params("name")
        return c.SendString("Hello, World with secrets 👋 !"+name+ " "+ secrets.GitHubAPIKey)
    })
}

// This is a simple REST API that responds with a personalized greeting.
// To call it, run in your terminal:
//
//     curl http://localhost:4000/hello/name
//
//encore:api public raw path=/hello/*p1
func Handle(w http.ResponseWriter, r *http.Request) {
	RemoveBaseURLFromRequest(r)
	
	// Call the app
	adaptor.FiberApp(app)(w, r)
}
