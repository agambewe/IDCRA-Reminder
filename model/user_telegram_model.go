package model

type UserTelegramModel struct {
	ID         string `gorm:"primaryKey;column:id"`
	IDTelegram int64  `gorm:"column:id_telegram;uniqueIndex"`
	FirstName  string `gorm:"column:first_name"`
	LastName   string `gorm:"column:last_name"`
	Username   string `gorm:"column:username"`
}

type Tabler interface {
	TableName() string
}

func (UserTelegramModel) TableName() string {
	return "users_telegram"
}
