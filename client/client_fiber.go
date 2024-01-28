package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func setUp(app *fiber.App, stage string) {
	fmt.Printf("Setting up %s stage.\n", stage)
	api := app.Group(fmt.Sprintf("/%s", stage))
	api.Post("/clients", createClient)
	api.Get("/clients/:id", getClientByID)
}

func createClient(c *fiber.Ctx) error {
	return c.Status(201).JSON(&fiber.Map{
		"message": "Client created",
	})
}

func getClientByID(c *fiber.Ctx) error {
	clientID := c.Params("id")
	return c.Status(200).JSON(fiber.Map{
		"message": fmt.Sprintf("Client ID: %s", clientID),
	})
}
