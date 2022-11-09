package service

import (
	"gorm.io/gorm"
	"idcra-telegram-scheduler/model"
)

type RecordServiceImpl struct{}

func NewRecordService() RecordService {
	return &RecordServiceImpl{}
}

func (r *RecordServiceImpl) RecordUserAnswer(db *gorm.DB, userRecord model.UserTelegramRecordModel) error {

	result := db.Create(&userRecord)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *RecordServiceImpl) IsAlreadyExist(db *gorm.DB, userRecord model.UserTelegramRecordModel) bool {

	var record model.UserTelegramRecordModel

	result := db.Where("DATE(created_at) = DATE(current_timestamp) AND answer_type = ?", userRecord.AnswerType).First(&record)
	if result.RowsAffected != 0 {
		return true
	}

	return false
}
