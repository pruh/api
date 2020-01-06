package messages

// OutboundTelegramMessage message to send to Telegram.
type OutboundTelegramMessage struct {
	ChatID              *int   `json:"chat_id"`
	DisablePreview      bool   `json:"disable_web_page_preview"`
	DisableNotification bool   `json:"disable_notification"`
	Text                string `json:"text"`
}

// InboundTelegramMessage message received by the server,
// that will be converted to OutboundTelegramMessage.
type InboundTelegramMessage struct {
	ChatID  *int   `json:"chat_id"`
	Message string `json:"message"`
	Silent  bool   `json:"silent"`
}

// NewOutboundTelegramMessage creates new OutboundTelegramMessage with default params.
func NewOutboundTelegramMessage(ChatID *int) OutboundTelegramMessage {
	return OutboundTelegramMessage{
		ChatID:              ChatID,
		DisablePreview:      true,
		DisableNotification: true,
	}
}

// NewInboundTelegramMessage creates new InboundTelegramMessage with default params.
func NewInboundTelegramMessage(ChatID *int) InboundTelegramMessage {
	return InboundTelegramMessage{
		ChatID: ChatID,
		Silent: true,
	}
}
