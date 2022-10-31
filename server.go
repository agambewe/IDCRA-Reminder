package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"idcra-telegram-scheduler/app"
	"idcra-telegram-scheduler/controller"
	"idcra-telegram-scheduler/helper"
)

func main() {
	err := godotenv.Load()
	helper.PanicIfError(err)

	db := app.NewDB()
	bot, err := tgbotapi.NewBotAPI(helper.Getenv("BOT_TOKEN", ""))
	helper.PanicIfError(err)

	botController := controller.NewBotController(bot, db)

	botController.SendDailyMessages()
	botController.ListenToBot()
}
