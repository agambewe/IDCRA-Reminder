package service

import (
	"gorm.io/gorm"
	"idcra-telegram-scheduler/model"
)

type RecordService interface {
	RecordUserAnswer(db *gorm.DB, userRecord model.UserTelegramRecordModel) error
	IsAlreadyExist(db *gorm.DB, userRecord model.UserTelegramRecordModel) bool
	CreateReport(db *gorm.DB) []model.UserRecordModel
}
