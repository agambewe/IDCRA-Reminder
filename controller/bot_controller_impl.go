package controller

import (
	"github.com/go-co-op/gocron"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
	"idcra-telegram-scheduler/helper"
	"log"
	"strconv"
	"time"
)

type BotControllerImpl struct {
	bot *tgbotapi.BotAPI
	db  *gorm.DB
}

func NewBotController(bot *tgbotapi.BotAPI, db *gorm.DB) BotController {
	return &BotControllerImpl{
		bot: bot,
		db:  db,
	}
}

func (b *BotControllerImpl) ListenToBot() {

	debugState := helper.Getenv("BOT_DEBUG", "false")
	debugStateBool, err := strconv.ParseBool(debugState)
	helper.PanicIfError(err)

	b.bot.Debug = debugStateBool
	log.Printf("Authorized on account %s", b.bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := b.bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // If we got a message
			continue
		}

		if !update.Message.IsCommand() { // ignore any non-command Messages
			continue
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
		switch update.Message.Command() {
		case "start":
			msg.Text = "Selamat anda sudah berlangganan bot kami, berikutnya anda akan mendapatkan daily reminder dari kami. Terimakasih."
		default:
			msg.Text = "Maaf, saat ini command tidak tersedia."
		}

		if _, err := b.bot.Send(msg); err != nil {
			log.Panic(err)
		}
	}
}

func (b *BotControllerImpl) StopListenToBot() {
	b.bot.StopReceivingUpdates()
}

func (b *BotControllerImpl) SendDailyMessages() {
	s := gocron.NewScheduler(time.UTC)

	_, err := s.Every(1).Day().At("06:00").At("21:00").Do(hitAuto, b.bot)
	helper.PanicIfError(err)

	s.StartImmediately()
	s.StartAsync()
}

func hitAuto(bot *tgbotapi.BotAPI) {
	msg := tgbotapi.NewMessage(1101320255, "Hallo")
	bot.Send(msg)
}
