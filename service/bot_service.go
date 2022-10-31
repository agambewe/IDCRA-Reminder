package service

import (
	"gorm.io/gorm"
	"idcra-telegram-scheduler/model"
)

type BotService interface {
	SaveUserTelegram(db *gorm.DB, userTelegram model.UserTelegramModel) error
	GetAllUserTelegram(db *gorm.DB) ([]model.UserTelegramModel, error)
}
