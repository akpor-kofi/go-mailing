package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
)

type jsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type handler struct {
	Mail Mail
}

func NewHandler(mail Mail) *handler {
	return &handler{Mail: mail}
}

func (h *handler) SendMail(c *fiber.Ctx) error {
	requestPayload := &struct {
		From    string `json:"from"`
		To      string `json:"to"`
		Subject string `json:"subject"`
		Message string `json:"message"`
	}{}

	if err := c.BodyParser(requestPayload); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	msg := Message{
		From:    requestPayload.From,
		To:      requestPayload.To,
		Subject: requestPayload.Subject,
		Data:    requestPayload.Message,
	}

	err := h.Mail.SendSMTPMessage(msg)
	fmt.Println(err)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	fmt.Println("no errors")

	payload := jsonResponse{
		Error:   false,
		Message: "sent to " + msg.To,
	}

	fmt.Println(payload)

	return c.Status(fiber.StatusAccepted).JSON(payload)
}
