package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/akpor-kofi/auth/models"
	"github.com/akpor-kofi/auth/ports"
	"github.com/gofiber/fiber/v2"
	"log"
	"net/http"
)

type jsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type FindUserByEmailer interface {
	GetByEmail(email string) (*models.User, error)
}

type authHandler struct {
	userStore FindUserByEmailer
}

func NewAuthHandler(userStore ports.UserRepository) *authHandler {
	return &authHandler{userStore: userStore}
}

func (h *authHandler) Authenticate(c *fiber.Ctx) error {
	log.Println("got here")
	requestPayload := &struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}

	if err := c.BodyParser(requestPayload); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	user, err := h.userStore.GetByEmail(requestPayload.Email)
	if err != nil {
		log.Println("invalid email")
		return fiber.NewError(fiber.StatusBadRequest, "invalid credentials")
	}

	if err := user.PasswordMatches(user.Password, requestPayload.Password); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid credentials")
	}

	// log authentication
	err = h.logRequest("authentication", fmt.Sprintf("%s logged in", user.Email))
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, err.Error())
	}

	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Logged in user %s", user.Email),
		Data:    user,
	}

	return c.Status(fiber.StatusAccepted).JSON(payload)
}

func (h *authHandler) logRequest(name, data string) error {
	entry := struct {
		Name string `json:"name"`
		Data string `json:"data"`
	}{
		Name: name,
		Data: data,
	}

	jsonData, _ := json.MarshalIndent(entry, "", "\t")
	logServiceURL := "http://logger-srv/log"

	_, err := http.Post(logServiceURL, "application/json", bytes.NewReader(jsonData))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return nil
}
