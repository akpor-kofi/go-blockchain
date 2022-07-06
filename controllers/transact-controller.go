package controllers

import (
	"fmt"
	"github.com/akpor-kofi/blockchain/app"
	"github.com/akpor-kofi/blockchain/components"
	"github.com/gofiber/fiber/v2"
)

func GenerateTransaction(ctx *fiber.Ctx) error {
	b := new(struct {
		Recipient string `json:"recipient"`
		Amount    int    `json:"amount"`
	})

	if err := ctx.BodyParser(b); err != nil {
		return err
	}

	transaction, ok := app.TransactionPool.ExistingTransaction(app.Wallet.PublicKey)
	var err error
	fmt.Println(transaction)
	if ok {
		fmt.Println("here")
		transaction.Update(app.Wallet, b.Recipient, b.Amount)
	} else {
		transaction, err = app.Wallet.CreateTransaction(b.Recipient, b.Amount, nil)
		if err != nil {
			return ctx.SendString(err.Error())
		}
	}

	app.TransactionPool.SetTransactions(transaction)
	app.BroadcastTransaction(transaction)

	return ctx.JSON(transaction)
}

func GetTransactionPoolMap(ctx *fiber.Ctx) error {

	return ctx.JSON(app.TransactionPool.TransactionMap)
}

func MineTransactions(ctx *fiber.Ctx) error {

	err := app.TransactMiner.MineTransaction()
	if err != nil {
		return err
	}

	return ctx.SendString("successfully mined transactions")
}

func GetWalletInfo(ctx *fiber.Ctx) error {

	return ctx.JSON(components.Wallet{
		PublicKey: app.Wallet.PublicKey,
		Balance:   components.CalculateBalance(app.BC.Chain, app.Wallet.PublicKey),
	})
}
