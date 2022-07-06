package controllers

import (
	"github.com/akpor-kofi/blockchain/app"
	"github.com/akpor-kofi/blockchain/components"
	"github.com/gofiber/fiber/v2"
)

func GetAllBlocks(ctx *fiber.Ctx) error {
	//fmt.Println(app.BC.Chain)

	ctx.SendStatus(200)
	return ctx.JSON(app.BC.Chain)
}

func AddBlock(ctx *fiber.Ctx) error {
	b := new(components.Block)

	if err := ctx.BodyParser(b); err != nil {
		return err
	}

	app.BC = app.BC.AddBlock(b.Data)

	app.BroadcastChain()

	return ctx.SendString("data saved")
}
