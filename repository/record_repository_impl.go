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

func (r *RecordRepositoryImpl) CreateReport(db *gorm.DB) map[int64]model.UserRecord {

	records := map[int64]model.UserRecord{}

	userRecords := r.recordService.CreateReport(db)

	for _, record := range userRecords {

		if _, ok := records[record.TelegramId]; !ok {
			records[record.TelegramId] = model.UserRecord{
				CountDayYES:   0,
				CountDayNO:    0,
				CountNightYES: 0,
				CountNightNO:  0,
			}
		}

		val, _ := records[record.TelegramId]

		if record.AnswerType == "DAY" {

			if record.UserAnswer == 1 {
				val.CountDayYES = record.AnsCount
			} else {
				val.CountDayNO = record.AnsCount
			}

		} else if record.AnswerType == "NIGHT" {

			if record.UserAnswer == 1 {
				val.CountNightYES = record.AnsCount
			} else {
				val.CountNightYES = record.AnsCount
			}
		}

		records[record.TelegramId] = val
	}

	return records
}
