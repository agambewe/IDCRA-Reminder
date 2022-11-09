package model

type Tabler interface {
	TableName() string
}

func (UserTelegramModel) TableName() string {
	return "users_telegram"
}

func (UserTelegramRecordModel) TableName() string {
	return "users_telegram_records"
}
