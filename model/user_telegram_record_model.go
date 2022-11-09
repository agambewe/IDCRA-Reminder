package model

type UserTelegramRecordModel struct {
	ID         string `gorm:"primaryKey;column:id"`
	AnswerType string `gorm:"column:answer_type"`
	UserAnswer bool   `gorm:"user_answer"`
	TelegramId int64  `gorm:"column:id_telegram"`
}

type UserRecordModel struct {
	AnswerType string `gorm:"column:answer_type"`
	TelegramId int64  `gorm:"column:id_telegram"`
	AnsCount   int    `gorm:"column:count"`
	UserAnswer int    `gorm:"column:user_answer"`
}
