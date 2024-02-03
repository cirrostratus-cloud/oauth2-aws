package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
)

func setUp(app *fiber.App, stage string) {
	log.
		WithField("Stage", stage).
		Info("Setting up stage.")
	api := app.Group(fmt.Sprintf("/%s", stage))
	api.Post("/users", createUser)
	api.Get("/users/:id", getUserByID)
}

func createUser(c *fiber.Ctx) error {
	return c.Status(201).JSON(&fiber.Map{
		"message": "User created",
	})
}

func getUserByID(c *fiber.Ctx) error {
	UserID := c.Params("id")
	return c.Status(200).JSON(fiber.Map{
		"message": fmt.Sprintf("User ID: %s", UserID),
	})
}
