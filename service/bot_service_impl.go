package service

import (
	"gorm.io/gorm"
	"idcra-telegram-scheduler/model"
)

type BotServiceImpl struct{}

func NewBotService() BotService {
	return &BotServiceImpl{}
}

func (b *BotServiceImpl) SaveUserTelegram(db *gorm.DB, userTelegram model.UserTelegramModel) error {
	result := db.Create(&userTelegram)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (b *BotServiceImpl) DeleteUserTelegram(db *gorm.DB, telegramID int64) error {

	result := db.Where("id_telegram = ?", telegramID).Delete(&model.UserTelegramModel{})
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (b *BotServiceImpl) GetAllUserTelegram(db *gorm.DB) ([]model.UserTelegramModel, error) {
	var usersData []model.UserTelegramModel

	result := db.Find(&usersData)
	if result.Error != nil {
		return []model.UserTelegramModel{}, result.Error
	}

	return usersData, nil
}

func (b *BotServiceImpl) GetUserByTelegramID(db *gorm.DB, telegramID int64) (bool, error) {

	result := db.First(&model.UserTelegramModel{}, "id_telegram = ?", telegramID)

	if result.RowsAffected == 0 {
		return false, result.Error
	}

	return true, result.Error
}
