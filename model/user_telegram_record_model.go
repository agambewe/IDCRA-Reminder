package model

type UserTelegramRecordModel struct {
	ID         string `gorm:"primaryKey;column:id"`
	UserAnswer bool   `gorm:"user_answer"`
	AnswerType string `gorm:"column:answer_type"`
	TelegramId int64  `gorm:"column:id_telegram"`
}
