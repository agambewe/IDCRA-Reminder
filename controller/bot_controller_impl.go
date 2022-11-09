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
	"reflect"
	"strconv"
	"time"
)

type MessageTime string

const (
	Morning MessageTime = "morning"
	Night               = "night"
)

type BotControllerImpl struct {
	bot              *tgbotapi.BotAPI
	db               *gorm.DB
	botRepository    repository.BotRepository
	recordRepository repository.RecordRepository
}

func NewBotController(bot *tgbotapi.BotAPI, db *gorm.DB, botRepository repository.BotRepository, recordRepository repository.RecordRepository) BotController {
	return &BotControllerImpl{
		bot:              bot,
		db:               db,
		botRepository:    botRepository,
		recordRepository: recordRepository,
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

		if update.Message != nil {

			if !update.Message.IsCommand() { // ignore any non-command Messages
				continue
			}

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
			switch update.Message.Command() {
			case "start", "subscribe":

				userTelegram := model.UserTelegram{
					IDTelegram: update.Message.Chat.ID,
					FirstName:  update.Message.Chat.FirstName,
					LastName:   update.Message.Chat.LastName,
					Username:   update.Message.Chat.UserName,
				}

				state, errDB := b.botRepository.SaveUserTelegram(b.db, userTelegram)
				if errDB != nil {
					sendMessageToDeveloper(b.bot, errDB.Error())
					msg.Text = "Maaf, terjadi kesalahan."
				} else {
					if !state {
						msg.Text = "Selamat anda sudah berlangganan bot kami, berikutnya anda akan mendapatkan daily reminder dari kami. Terimakasih."
					} else {
						msg.Text = "Anda sudah berlangganan bot kami, harap menunggu daily reminder dari kami. Terimakasih."
					}
				}
			case "unsubscribe":
				state, _ := b.botRepository.DeleteUserTelegram(b.db, update.Message.Chat.ID)

				if errDB != nil && errDB != gorm.ErrRecordNotFound {

					sendMessageToDeveloper(b.bot, errDB.Error())
					msg.Text = "Maaf terjadi kesalahan."
				} else {
					if state {
						msg.Text = "Terimakasih telah berlangganan bot kami, kami akan menghentikan layanan daily reminder."
					} else {
						msg.Text = "Maaf, anda belum berlangganan bot kami sebelumnya."
					}
				}

			default:
				msg.Text = "Maaf, saat ini command tidak tersedia."
			}

			if _, err := b.bot.Send(msg); err != nil {
				helper.PanicIfError(err)
			}

			helper.PanicIfError(errDB)

		} else if update.CallbackQuery != nil {
			//1 == YES MORNING
			//2 == NO MORNING
			//3 == YES NIGHT
			//4 == NO NIGHT
			message := ""
			request := model.UserTelegramRecord{}

			switch update.CallbackQuery.Data {
			case "1":

				request = model.UserTelegramRecord{
					UserAnswer: "SUDAH",
					AnswerType: "DAY",
					TelegramId: update.CallbackQuery.Message.Chat.ID,
				}

				message = "Dengan sikat gigi setelah sarapan di pagi hari kamu dapat mencegah timbulnya gigi berlubang! 😃"

			case "2":

				request = model.UserTelegramRecord{
					UserAnswer: "BELUM",
					AnswerType: "DAY",
					TelegramId: update.CallbackQuery.Message.Chat.ID,
				}

				message = "Segera sikat gigi sebelum kuman - kuman dalam mulut membuat gigi mu berlubang dan akan muncul rasa sakit! \U0001FAE3"

			case "3":
				request = model.UserTelegramRecord{
					UserAnswer: "SUDAH",
					AnswerType: "NIGHT",
					TelegramId: update.CallbackQuery.Message.Chat.ID,
				}

				message = "Dengan sikat gigi sebelum tidur di malam hari kamu dapat mencegah timbulnya gigi berlubang! 😃"

			case "4":
				request = model.UserTelegramRecord{
					UserAnswer: "BELUM",
					AnswerType: "NIGHT",
					TelegramId: update.CallbackQuery.Message.Chat.ID,
				}

				message = "Segera sikat gigi sebelum kuman - kuman dalam mulut membuat gigi mu berlubang dan akan muncul rasa sakit! \U0001FAE3"

			}

			state, err := b.recordRepository.RecordUserAnswer(b.db, request)

			if state {
				respondToInlineInput(
					b.bot,
					update.CallbackQuery,
					message,
				)
			} else {

				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "")
				if err != nil {
					sendMessageToDeveloper(b.bot, errDB.Error())
					msg.Text = "Maaf terjadi kesalahan."
				} else {
					msg.Text = "Anda sudah memberikan jawaban."
				}

				if _, err := b.bot.Send(msg); err != nil {
					helper.PanicIfError(err)
				}
			}
		}
	}
}

func (b *BotControllerImpl) StopListenToBot() {
	b.bot.StopReceivingUpdates()
}

func (b *BotControllerImpl) SendDailyMessages() {
	s := gocron.NewScheduler(time.UTC)

	_, err := s.Every(1).Day().At("06:00").Tag("morning").Do(sendMessages, b.bot, b.db, b.botRepository, reflect.ValueOf(Morning).String())
	helper.PanicIfError(err)

	_, err = s.Every(1).Day().At("21:00").Tag("night").Do(sendMessages, b.bot, b.db, b.botRepository, reflect.ValueOf(Night).String())
	helper.PanicIfError(err)

	_, err = s.Cron("0 7 1 1/1 *").Tag("report").Do(sendReports, b.bot, b.db, b.recordRepository)
	helper.PanicIfError(err)

	//s.StartImmediately()
	s.StartAsync()
	s.RunAll()

	sendMessageToDeveloper(b.bot, "Scheduler Running...")
}

func sendReports(bot *tgbotapi.BotAPI, db *gorm.DB, recordRepository repository.RecordRepository) {
	userRecords := recordRepository.CreateReport(db)
	log.Print(userRecords)
	for key, val := range userRecords {
		msg := tgbotapi.NewMessage(key, fmt.Sprintf(`Dalam sebulan ini sikat gigi pagi setelah sarapan %v kali, sikat gigi malam sebelum tidur %v kali, dan telah terlewat %v kali.

Ayo jaga terus kesehatan gigi mu dengan sikat gigi 2x sehari pagi setelah sarapan dan malam sebelum tidur
		`, val.CountDayYES, val.CountNightYES, val.CountDayNO+val.CountNightNO))

		if _, err := bot.Send(msg); err != nil {
			panic(err)
		}
	}
}

func sendMessages(bot *tgbotapi.BotAPI, db *gorm.DB, botRepository repository.BotRepository, when string) {

	usersTelegram := botRepository.GetAllUsersTelegram(db)

	for _, user := range usersTelegram {
		if when == reflect.ValueOf(Morning).String() {

			numericKeyboard := tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("Sudah", "1"),
					tgbotapi.NewInlineKeyboardButtonData("Belum", "2"),
				),
			)

			sendMessageToUserWithKeyboard(bot, &numericKeyboard, user.IDTelegram, fmt.Sprintf("Halo %s %s, sudahkah kamu sikat gigi setelah sarapan pagi ini? 😊", user.FirstName, user.LastName))

		} else if when == reflect.ValueOf(Night).String() {

			numericKeyboard := tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("Sudah", "3"),
					tgbotapi.NewInlineKeyboardButtonData("Belum", "4"),
				),
			)

			sendMessageToUserWithKeyboard(bot, &numericKeyboard, user.IDTelegram, fmt.Sprintf("Halo %s %s, sudahkah kamu sikat gigi sebelum tidur malam ini? 😊", user.FirstName, user.LastName))
		}
	}
}

func sendMessageToUserWithKeyboard(bot *tgbotapi.BotAPI, keyboard *tgbotapi.InlineKeyboardMarkup, telegramId int64, msgInput string) {

	msg := tgbotapi.NewMessage(telegramId, msgInput)
	msg.ReplyMarkup = *keyboard

	_, err := bot.Send(msg)
	if err != nil {
		log.Printf("Error when sending message: %v", err)
	}
}

// - Alert Developer if Error happened
func sendMessageToDeveloper(bot *tgbotapi.BotAPI, msgInput string) {

	devIDString := helper.Getenv("DEV_CHAT_ID", "")
	devIDInt, _ := strconv.Atoi(devIDString)

	msg := tgbotapi.NewMessage(int64(devIDInt), msgInput)

	_, err := bot.Send(msg)
	if err != nil {
		log.Printf("Error when sending message: %v", err)
	}
}

func respondToInlineInput(bot *tgbotapi.BotAPI, callbackQuery *tgbotapi.CallbackQuery, message string) {
	//callback := tgbotapi.NewCallback(callbackQuery.ID, callbackQuery.Data)
	//if _, err := bot.Request(callback); err != nil {
	//	panic(err)
	//}

	msg := tgbotapi.NewMessage(callbackQuery.Message.Chat.ID, message)
	if _, err := bot.Send(msg); err != nil {
		panic(err)
	}
}
