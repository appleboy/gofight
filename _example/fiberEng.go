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
	e.Get("/helloCouple", fiberRoute)
	return e
}

func fiberRoute(c *fiber.Ctx) error {
	names := c.Context().QueryArgs().PeekMulti("names")
	ages := c.Context().QueryArgs().PeekMulti("ages")
	msg := ""
	for i := range names {
		msg += fmt.Sprintf("God Love the World ! 👴 %s is %s years old~\n", string(names[i]), string(ages[i]))
	}
	c.Status(200).SendString(msg)
	// =>God Love the World ! 👴 john is 75 years old~
	//   God Love the World ! 👴 mary is 25 years old~
	return nil
}
