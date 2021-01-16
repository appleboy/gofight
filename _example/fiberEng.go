package example

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
)

// FiberEngine is fiber web server and handlers routes
func FiberEngine() *fiber.App {
	// Fiber instance
	e := fiber.New()
	// Routes
	e.Get("/", fiberRoute)
	return e
}

func fiberRoute(c *fiber.Ctx) error {
	msg := fmt.Sprintf("God Love the World ! 👴 %s is %s years old~", c.Query("name"), c.Query("age"))
	c.Status(200).SendString(msg) // =>God Love the World ! 👴 john is 75 years old~
	return nil
}
