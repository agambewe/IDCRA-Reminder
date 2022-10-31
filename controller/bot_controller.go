package controller

type BotController interface {
	ListenToBot()
	StopListenToBot()
	SendDailyMessages()
}
