package sendpulse_sdk_go

type BotsService struct {
	client   *Client
	Fb       *BotsFbService
	Vk       *BotsVkService
	Telegram *BotsTelegramService
	WhatsApp *BotsWhatsAppService
	Ig       *BotsIgService
	LiveChat *BotsLiveChatService
}

func newBotsService(cl *Client) *BotsService {
	return &BotsService{
		client:   cl,
		Fb:       newBotsFbService(cl),
		Vk:       newBotsVkService(cl),
		Telegram: newBotsTelegramService(cl),
		WhatsApp: newBotsWhatsAppService(cl),
		Ig:       newBotsIgService(cl),
		LiveChat: newBotsLiveChatService(cl),
	}
}
