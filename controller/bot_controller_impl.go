package controller

import (
	"fmt"
	"github.com/go-co-op/gocron"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
	"idcra-telegram-scheduler/helper"
	"idcra-telegram-scheduler/model"
	"idcra-telegram-scheduler/repository"
	"log"
	"strconv"
	"time"
)

type BotControllerImpl struct {
	bot           *tgbotapi.BotAPI
	db            *gorm.DB
	botRepository repository.BotRepository
}

func NewBotController(bot *tgbotapi.BotAPI, db *gorm.DB, botRepository repository.BotRepository) BotController {
	return &BotControllerImpl{
		bot:           bot,
		db:            db,
		botRepository: botRepository,
	}
}

func (b *BotControllerImpl) ListenToBot() {

	var errDB error

	debugState := helper.Getenv("BOT_DEBUG", "false")
	debugStateBool, err := strconv.ParseBool(debugState)
	helper.PanicIfError(err)

	b.bot.Debug = debugStateBool
	log.Printf("Authorized on account %s", b.bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := b.bot.GetUpdatesChan(u)

	defer helper.RecoveryIfPanic(b.db)

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

			userTelegram := model.UserTelegram{
				IDTelegram: update.Message.Chat.ID,
				FirstName:  update.Message.Chat.FirstName,
				LastName:   update.Message.Chat.LastName,
				Username:   update.Message.Chat.UserName,
			}

			errDB = b.botRepository.SaveUserTelegram(b.db, userTelegram)
			if errDB != nil {

				sendMessageToDeveloper(b.bot, errDB.Error())
				msg.Text = "Maaf terjadi kesalahan."
			} else {
				msg.Text = "Selamat anda sudah berlangganan bot kami, berikutnya anda akan mendapatkan daily reminder dari kami. Terimakasih."
			}

		default:
			msg.Text = "Maaf, saat ini command tidak tersedia."
		}

		if _, err := b.bot.Send(msg); err != nil {
			helper.PanicIfError(err)
		}

		helper.PanicIfError(errDB)
	}
}

func (b *BotControllerImpl) StopListenToBot() {
	b.bot.StopReceivingUpdates()
}

func (b *BotControllerImpl) SendDailyMessages() {
	s := gocron.NewScheduler(time.UTC)

	_, err := s.Every(1).Day().At("06:00").At("21:00").Do(sendMessages, b.bot, b.db, b.botRepository)
	helper.PanicIfError(err)

	s.StartImmediately()
	s.StartAsync()

	sendMessageToDeveloper(b.bot, "Scheduler Running...")
}

// - Alert Developer if Error happened
func sendMessageToDeveloper(bot *tgbotapi.BotAPI, msgInput string) {

	devIDString := helper.Getenv("DEV_CHAT_ID", "")
	devIDInt, _ := strconv.Atoi(devIDString)

	msg := tgbotapi.NewMessage(int64(devIDInt), msgInput)
	bot.Send(msg)
}

func sendMessages(bot *tgbotapi.BotAPI, db *gorm.DB, botRepository repository.BotRepository) {

	usersTelegram := botRepository.GetAllUsersTelegram(db)

	for _, user := range usersTelegram {

		msg := tgbotapi.NewMessage(
			user.IDTelegram,
			fmt.Sprintf("Hai %s %s, jangan lupa membersihkan gigi ya", user.FirstName, user.LastName),
		)
		bot.Send(msg)
	}
}
