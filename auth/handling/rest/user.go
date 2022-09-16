package rest

import (
	"github.com/akpor-kofi/auth/models"
	"github.com/akpor-kofi/auth/ports"
	"github.com/gofiber/fiber/v2"
)

type userHandler struct {
	userStore ports.UserRepository
}

func NewUserHandler(userStore ports.UserRepository) *userHandler {
	return &userHandler{userStore}
}

func (h *userHandler) CreateUser(c *fiber.Ctx) error {
	user := new(models.User)

	if err := c.BodyParser(user); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if err := h.userStore.Create(user); err != nil {

		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(user)
}

func (h *userHandler) GetAllUsers(c *fiber.Ctx) error {
	list, err := h.userStore.List()
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(list)
}

func (h *userHandler) GetUser(c *fiber.Ctx) error {
	userId := c.Params("id")

	user, err := h.userStore.Get(userId)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(user)
}

func (h *userHandler) UpdateUser(c *fiber.Ctx) error {
	userId := c.Params("id")

	user := new(models.User)

	if err := c.BodyParser(user); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	updatedUser, err := h.userStore.Update(userId, user)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(updatedUser)
}
