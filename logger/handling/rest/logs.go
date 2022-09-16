package rest

import (
	"github.com/akpor-kofi/logger/models"
	"github.com/akpor-kofi/logger/ports"
	"github.com/gofiber/fiber/v2"
	"log"
)

type jsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type JSONPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

type logEntryHandler struct {
	logStore ports.LogEntryRepository
}

func NewLogEntryHandler(logStore ports.LogEntryRepository) *logEntryHandler {
	return &logEntryHandler{logStore: logStore}
}

func (leh *logEntryHandler) WriteLog(c *fiber.Ctx) error {
	requestPayload := new(JSONPayload)

	if err := c.BodyParser(requestPayload); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	log.Println(*requestPayload)

	event := models.LogEntry{
		Name: requestPayload.Name,
		Data: requestPayload.Data,
	}

	err := leh.logStore.Insert(&event)
	if err != nil {
		log.Println(err.Error())
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	log.Println("here")

	resp := jsonResponse{
		Error:   false,
		Message: "logged",
	}

	return c.Status(fiber.StatusCreated).JSON(resp)
}
