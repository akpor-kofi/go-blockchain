package routes

import (
	"github.com/akpor-kofi/blockchain/controllers"
	"github.com/gofiber/fiber/v2"
)

func BlockRouter(router fiber.Router) {
	router.Get("/blocks", controllers.GetAllBlocks)

	router.Post("/mine", controllers.AddBlock)
}
