package messages

// TelegramMessage message to send to Telegram
type TelegramMessage struct {
	ChatID              *int   `json:"chat_id"`
	DisablePreview      bool   `json:"disable_web_page_preview"`
	DisableNotification bool   `json:"disable_notification"`
	Text                string `json:"text"`
}

// Message message received by the server
type Message struct {
	ChatID  *int   `json:"chat_id"`
	Message string `json:"message"`
	Silent  bool   `json:"silent"`
}

// NewTelegramMessage creates new TelegramMessage with default params
func NewTelegramMessage(ChatID *int) TelegramMessage {
	return TelegramMessage{
		ChatID:              ChatID,
		DisablePreview:      true,
		DisableNotification: true,
	}
}

// NewMessage creates new Message with default params.
func NewMessage(ChatID *int) Message {
	return Message{
		ChatID: ChatID,
		Silent: true,
	}
}
