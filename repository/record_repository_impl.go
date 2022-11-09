package repository

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"idcra-telegram-scheduler/model"
	"idcra-telegram-scheduler/service"
)

type RecordRepositoryImpl struct {
	recordService service.RecordService
}

func NewRecordRepository(recordService service.RecordService) RecordRepository {
	return &RecordRepositoryImpl{recordService: recordService}
}

func (r *RecordRepositoryImpl) RecordUserAnswer(db *gorm.DB, request model.UserTelegramRecord) (bool, error) {

	newId := uuid.NewV1()
	newIdString := newId.String()

	var ans bool

	if request.UserAnswer == "SUDAH" {
		ans = true
	} else if request.UserAnswer == "BELUM" {
		ans = false
	}

	userRecord := model.UserTelegramRecordModel{
		ID:         newIdString,
		UserAnswer: ans,
		AnswerType: request.AnswerType,
		TelegramId: request.TelegramId,
	}

	state := r.recordService.IsAlreadyExist(db, userRecord)

	if !state {
		err := r.recordService.RecordUserAnswer(db, userRecord)
		if err != nil {
			return false, err
		}
	} else {
		return false, nil
	}

	return true, nil
}
