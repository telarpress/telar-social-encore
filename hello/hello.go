package hello

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
)

// Cache state
var app *fiber.App



// init
func init() {



	// Initialize app
	app = fiber.New()
	app.Get("/:name", func(c *fiber.Ctx) error {
		fmt.Println("[ "+os.Getenv("PORT")+" ]")
		name := c.Params("name")
        return c.SendString("Hello, World with secrets 👋 !"+name )
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
