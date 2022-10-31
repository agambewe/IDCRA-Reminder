package repository

import (
	"gorm.io/gorm"
	"idcra-telegram-scheduler/model"
)

type BotRepository interface {
	SaveUserTelegram(db *gorm.DB, request model.UserTelegram) error
	GetAllUsersTelegram(db *gorm.DB) []model.UserTelegram
}