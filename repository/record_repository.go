package repository

import (
	"gorm.io/gorm"
	"idcra-telegram-scheduler/model"
)

type RecordRepository interface {
	RecordUserAnswer(db *gorm.DB, request model.UserTelegramRecord) (bool, error)
	CreateReport(db *gorm.DB) map[int64]model.UserRecord
}
