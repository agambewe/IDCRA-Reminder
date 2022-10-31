package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"idcra-telegram-scheduler/app"
	"idcra-telegram-scheduler/controller"
	"idcra-telegram-scheduler/helper"
	"idcra-telegram-scheduler/repository"
	"idcra-telegram-scheduler/service"
)

func main() {
	err := godotenv.Load()
	helper.PanicIfError(err)

	db := app.NewDB()
	botService := service.NewBotService()
	botRepository := repository.NewBotRepository(botService)

	bot, err := tgbotapi.NewBotAPI(helper.Getenv("BOT_TOKEN", ""))
	helper.PanicIfError(err)

	botController := controller.NewBotController(bot, db, botRepository)

	botController.SendDailyMessages()
	botController.ListenToBot()
}
