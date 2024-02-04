package main

import (
	"fmt"

	"github.com/cirrostratus-cloud/oauth2/user"
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
)

type userRequest struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type userResponse struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Enabled   bool   `json:"enabled"`
}

type userAPI struct {
	api               fiber.Router
	createUserUseCase user.CreateUserUseCase
	getUserUseCase    user.GetUserUseCase
}

func newUserAPI(createUserUseCase user.CreateUserUseCase, getUserUseCase user.GetUserUseCase) *userAPI {
	return &userAPI{
		createUserUseCase: createUserUseCase,
		getUserUseCase:    getUserUseCase,
	}
}

func (u *userAPI) setUp(app *fiber.App, stage string) (userAPI *userAPI) {
	log.
		WithField("Stage", stage).
		Info("Setting up stage.")
	u.api = app.Group(fmt.Sprintf("/%s", stage))
	u.api.Post("/users", u.createUser)
	u.api.Get("/users/:id", u.getUserByID)
	return u
}

func (u *userAPI) createUser(c *fiber.Ctx) error {
	userRequest := new(userRequest)
	if err := c.BodyParser(userRequest); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   "Bad Request",
			"message": err.Error(),
		})
	}
	user, err := u.createUserUseCase.NewUser(user.CreateUserRequest{
		Email:     userRequest.Email,
		Password:  userRequest.Password,
		FirstName: userRequest.FirstName,
		LastName:  userRequest.LastName,
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Internal Server Error",
			"message": err.Error(),
		})
	}
	return c.Status(201).JSON(&fiber.Map{
		"user": userResponse{
			ID: user.UserID,
		},
		"message": "User created",
	})
}

func (u *userAPI) getUserByID(c *fiber.Ctx) error {
	UserID := c.Params("id")
	user, err := u.getUserUseCase.GetUserByID(user.UserByID{UserID: UserID})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Internal Server Error",
			"message": err.Error(),
		})
	}
	return c.Status(200).JSON(fiber.Map{
		"user": userResponse{
			ID:        user.UserID,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Enabled:   user.Enabled,
		},
		"message": fmt.Sprintf("User ID: %s", UserID),
	})
}
