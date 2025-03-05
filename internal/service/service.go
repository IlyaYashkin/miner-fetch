package service

type Service struct {
	Device         *DeviceService
	Polling        *Polling
	Logger         *Logger
	TelegramSender *TelegramSender
	HttpClient     *HttpClient
}
