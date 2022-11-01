package repository

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"idcra-telegram-scheduler/helper"
	"idcra-telegram-scheduler/model"
	"idcra-telegram-scheduler/service"
)

type BotRepositoryImpl struct {
	botService service.BotService
}

func NewBotRepository(botService service.BotService) BotRepository {
	return &BotRepositoryImpl{
		botService: botService,
	}
}

func (b *BotRepositoryImpl) SaveUserTelegram(db *gorm.DB, request model.UserTelegram) (bool, error) {

	state, err := b.botService.GetUserByTelegramID(db, request.IDTelegram)

	if !state {
		newId := uuid.NewV1()
		newIdString := newId.String()

		userTelegram := model.UserTelegramModel{
			ID:         newIdString,
			IDTelegram: request.IDTelegram,
			FirstName:  request.FirstName,
			LastName:   request.LastName,
			Username:   request.Username,
		}

		err = b.botService.SaveUserTelegram(db, userTelegram)
	}

	return state, err
}

func (b *BotRepositoryImpl) GetAllUsersTelegram(db *gorm.DB) []model.UserTelegram {

	defer func() {
		helper.RecoveryIfPanic(db)
	}()

	usersData, err := b.botService.GetAllUserTelegram(db)
	helper.PanicIfError(err)

	var users []model.UserTelegram

	for _, user := range usersData {
		userTelegram := model.UserTelegram{
			IDTelegram: user.IDTelegram,
			FirstName:  user.FirstName,
			LastName:   user.LastName,
			Username:   user.Username,
		}

		users = append(users, userTelegram)
	}

	return users
}
