package routes

import (
	"github.com/akpor-kofi/blockchain/controllers"
	"github.com/gofiber/fiber/v2"
)

func TransactRouter(router fiber.Router) {
	router.Get("/transaction-pool-map", controllers.GetTransactionPoolMap)
	router.Post("/transact", controllers.GenerateTransaction)
	router.Get("/mine-transactions", controllers.MineTransactions)
	router.Get("/wallet-info", controllers.GetWalletInfo)
}
