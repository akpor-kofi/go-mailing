package main

import (
	"fmt"
	"github.com/akpor-kofi/broker/handling/rest"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"math"
	"time"
)

const webPort = ":8080"

func main() {
	rabbitConn, err := connect()
	if err != nil {
		log.Fatal(err)
	}

	defer rabbitConn.Close()

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "https://*, http://*",
		AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS",
		AllowHeaders:     "Accept, Authorization, Content-Type, X-CSRF-Token",
		ExposeHeaders:    "Link",
		AllowCredentials: true,
		MaxAge:           300,
	}))

	brokerHandler := rest.NewBroker(rabbitConn)

	app.Post("/", brokerHandler.Broker)

	app.Post("/handle", brokerHandler.HandleSubmission)

	log.Fatal(app.Listen(webPort))
}

func connect() (*amqp.Connection, error) {
	var counts int64
	var backOff = 1 * time.Second
	var connection *amqp.Connection

	// don't continue until rabbit is ready

	for {
		c, err := amqp.Dial("amqp://guest:guest@rabbitmq")
		if err != nil {
			fmt.Println("rabbitMQ not yet ready")
			counts++
		} else {
			connection = c
			break
		}

		if counts > 5 {
			fmt.Println(err)
			return nil, err
		}

		backOff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		time.Sleep(backOff)
		continue
	}

	return connection, nil
}
