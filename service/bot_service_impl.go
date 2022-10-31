package service

import (
	"gorm.io/gorm"
	"idcra-telegram-scheduler/model"
)

type BotServiceImpl struct{}

func NewBotService() BotService {
	return &BotServiceImpl{}
}

func (b BotServiceImpl) SaveUserTelegram(db *gorm.DB, userTelegram model.UserTelegramModel) error {
	result := db.Create(&userTelegram)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (b BotServiceImpl) GetAllUserTelegram(db *gorm.DB) ([]model.UserTelegramModel, error) {
	var usersData []model.UserTelegramModel

	result := db.Find(&usersData)
	if result.Error != nil {
		return []model.UserTelegramModel{}, result.Error
	}

	return usersData, nil
}
