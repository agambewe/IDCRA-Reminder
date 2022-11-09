package model

type UserTelegramRecord struct {
	UserAnswer string
	AnswerType string
	TelegramId int64
}

type UserRecord struct {
	CountDayYES   int
	CountDayNO    int
	CountNightYES int
	CountNightNO  int
}
