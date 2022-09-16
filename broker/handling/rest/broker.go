package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/akpor-kofi/broker/event"
	"github.com/gofiber/fiber/v2"
	amqp "github.com/rabbitmq/amqp091-go"
	"net/http"
	"net/rpc"
)

type RPCPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

type jsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
	Log    LogPayload  `json:"log,omitempty"`
	Mail   MailPayload `json:"mail,omitempty"`
}

type MailPayload struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

type BrokerHandler struct {
	conn *amqp.Connection
}

func NewBroker(conn *amqp.Connection) *BrokerHandler {
	return &BrokerHandler{
		conn: conn,
	}
}

func (bh *BrokerHandler) Broker(c *fiber.Ctx) error {
	payload := jsonResponse{
		Error:   false,
		Message: "hit the broker",
	}

	return c.Status(fiber.StatusAccepted).JSON(payload)
}

func (bh *BrokerHandler) HandleSubmission(c *fiber.Ctx) error {
	requestPayload := new(RequestPayload)

	if err := c.BodyParser(requestPayload); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	switch requestPayload.Action {
	case "auth":
		return bh.authenticate(c, requestPayload.Auth)
	case "log":
		return bh.logItemViaRPC(c, requestPayload.Log)
	case "mail":
		return bh.sendMail(c, requestPayload.Mail)

	default:
		return fiber.NewError(fiber.StatusBadRequest, "unknown action")
	}
}

func (bh *BrokerHandler) sendMail(c *fiber.Ctx, msg MailPayload) error {
	jsonData, _ := json.MarshalIndent(msg, "", "\t")

	mailServiceURL := "http://mailer-srv/send"

	resp, err := http.Post(mailServiceURL, "application/json", bytes.NewReader(jsonData))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	defer resp.Body.Close()

	fmt.Println(resp.StatusCode)

	if resp.StatusCode != fiber.StatusAccepted {
		return fiber.NewError(fiber.StatusInternalServerError, "error calling mail service")
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "Message sent to " + msg.To

	return c.Status(fiber.StatusAccepted).JSON(payload)
}

func (bh *BrokerHandler) logItem(c *fiber.Ctx, l LogPayload) error {
	jsonData, _ := json.MarshalIndent(l, "", "\t")

	logServiceURL := "http://logger-srv/log"

	resp, err := http.Post(logServiceURL, "application/json", bytes.NewReader(jsonData))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusCreated {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "logged!"

	return c.Status(fiber.StatusAccepted).JSON(payload)

}

func (bh *BrokerHandler) authenticate(c *fiber.Ctx, a AuthPayload) error {
	jsonData, _ := json.MarshalIndent(a, "", "\t")

	resp, err := http.Post("http://auth-srv/authenticate", "application/json", bytes.NewReader(jsonData))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	defer resp.Body.Close()

	if resp.StatusCode == fiber.StatusUnauthorized {
		return fiber.NewError(fiber.StatusBadRequest, "invlaid credentials")
	} else if resp.StatusCode != fiber.StatusAccepted {
		return fiber.NewError(fiber.StatusInternalServerError, "error calling auth service")
	}

	jsonFromService := new(jsonResponse)

	err = json.NewDecoder(resp.Body).Decode(jsonFromService)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if jsonFromService.Error {
		return fiber.NewError(fiber.StatusUnauthorized, err.Error())
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "Authenticated!"
	payload.Data = jsonFromService.Data

	return c.Status(fiber.StatusAccepted).JSON(payload)
}

func (bh *BrokerHandler) logEventViaRabbit(c *fiber.Ctx, l LogPayload) error {
	err := bh.pushToQueue(l.Name, l.Data)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "logged via rabbitMQ"

	return c.Status(fiber.StatusAccepted).JSON(payload)
}

func (bh *BrokerHandler) logItemViaRPC(c *fiber.Ctx, l LogPayload) error {
	client, err := rpc.Dial("tcp", "logger-srv:5001")
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	rpcPayload := RPCPayload{
		Name: l.Name,
		Data: l.Data,
	}

	var result string
	err = client.Call("RPCServer.LogInfo", rpcPayload, &result)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	payload := jsonResponse{
		Error:   false,
		Message: result,
	}

	return c.Status(fiber.StatusAccepted).JSON(payload)
}

func (bh *BrokerHandler) pushToQueue(name, msg string) error {
	emitter, err := event.NewEventEmitter(bh.conn)
	if err != nil {
		return err
	}

	payload := LogPayload{Name: name, Data: msg}

	j, _ := json.MarshalIndent(&payload, "", "\t")
	err = emitter.Push(string(j), "log.INFO")
	if err != nil {
		return err
	}

	return nil
}
