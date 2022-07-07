package main

import (
	"fmt"
	pubsub "github.com/akpor-kofi/blockchain/app"
	"github.com/akpor-kofi/blockchain/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"log"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"
)

var Wg = sync.WaitGroup{}

func main() {
	app := fiber.New()
	api := app.Group("/api")
	api.Route("/", routes.BlockRouter)
	api.Route("/", routes.TransactRouter)

	Wg.Add(1)

	go pubsub.ConsumeEvents()

	addr := getAddressPort()

	err := app.Listen(addr)
	if err != nil {
		panic(err.Error())
	}

	Wg.Wait()
}

func getEnvVariable(key string) (string, error) {

	err := godotenv.Load("app.env")
	if err != nil {
		return "", err
	}

	return os.Getenv(key), nil
}

func getAddressPort() string {
	defaultPort := 5000
	var peerPort int
	var port int

	herokuPort := os.Getenv("PORT")
	fmt.Println(herokuPort, "here")
	if herokuPort != "" {
		return ":" + strconv.FormatInt(int64(port), 10)
	}

	genPeerPort, err := getEnvVariable("GENERATE_PEER_PORT")
	if err != nil {
		log.Fatal(err)
	}
	if genPeerPort == "true" {
		rand.Seed(time.Now().UnixNano()) // if this is not setup, it results in a deterministic random value which would cause conflicting ports
		peerPort = defaultPort + rand.Intn(100)
	}

	if peerPort > 5000 {
		port = peerPort
	} else {
		port = defaultPort
	}

	addr := ":" + strconv.FormatInt(int64(port), 10)

	return addr
}
