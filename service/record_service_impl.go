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

	result := db.Where("DATE(created_at) = DATE(current_timestamp) AND answer_type = ? AND id_telegram = ?", userRecord.AnswerType, userRecord.TelegramId).First(&record)
	if result.RowsAffected != 0 {
		return true
	}

	return false
}

func (r *RecordServiceImpl) CreateReport(db *gorm.DB) []model.UserRecordModel {

	var record []model.UserRecordModel

	result := db.Raw(`select id_telegram, answer_type, user_answer, count(user_answer) as count from users_telegram_records
                                                             where datediff(current_date, date(created_at)) <= 30
                                                             group by id_telegram, answer_type, user_answer
                                                             order by id_telegram, answer_type`).Find(&record)

	if result.RowsAffected == 0 {
		return []model.UserRecordModel{}
	}

	return record
}
