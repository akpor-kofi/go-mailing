package main

import (
	"fmt"
	"github.com/akpor-kofi/logger/db/logging"
	"github.com/akpor-kofi/logger/handling/rest"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"log"
	"net"
	"net/rpc"
)

const (
	webPort  = ":80"
	rpcPort  = "5001"
	mongoURL = "mongodb://logging:27017"
	gRpcPort = "50001"
)

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	log.Println("Starting logger service...")

	logStore := logging.NewLogStore()
	logHandler := rest.NewLogEntryHandler(logStore)

	err := rpc.Register(NewRPCServer(logStore))

	if err != nil {
		log.Fatalln(err)
	}
	go rpcListen()

	app := fiber.New(fiber.Config{
		ErrorHandler: rest.GlobalHandler,
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "https://*, http://*",
		AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS",
		AllowHeaders:     "Accept, Authorization, Content-Type, X-CSRF-Token",
		ExposeHeaders:    "Link",
		AllowCredentials: true,
		MaxAge:           300,
	}))

	app.Post("/log", logHandler.WriteLog)

	log.Fatal(app.Listen(webPort))
}

func rpcListen() {
	log.Println("Starting rpc listener service...5001")
	listen, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", rpcPort))
	if err != nil {
		panic(err)
	}
	defer listen.Close()

	for {
		rpcConn, err := listen.Accept()
		if err != nil {
			continue
		}

		go rpc.ServeConn(rpcConn)
	}
}
