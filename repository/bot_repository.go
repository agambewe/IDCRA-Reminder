package repository

import (
	"gorm.io/gorm"
	"idcra-telegram-scheduler/model"
)

type BotRepository interface {
	SaveUserTelegram(db *gorm.DB, request model.UserTelegram) (bool, error)
	DeleteUserTelegram(db *gorm.DB, idTelegram int64) (bool, error)
	GetAllUsersTelegram(db *gorm.DB) []model.UserTelegram
}
