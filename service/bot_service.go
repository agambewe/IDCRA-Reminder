package service

import (
	"gorm.io/gorm"
	"idcra-telegram-scheduler/model"
)

type BotService interface {
	SaveUserTelegram(db *gorm.DB, userTelegram model.UserTelegramModel) error
	DeleteUserTelegram(db *gorm.DB, telegramID int64) error
	GetAllUserTelegram(db *gorm.DB) ([]model.UserTelegramModel, error)
	GetUserByTelegramID(db *gorm.DB, telegramID int64) (bool, error)
}
