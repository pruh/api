package models

type OutboundTelegramMessage struct {
	ChatID              *int   `json:"chat_id"`
	DisablePreview      bool   `json:"disable_web_page_preview"`
	DisableNotification bool   `json:"disable_notification"`
	Text                string `json:"text"`
}

type InboundTelegramMessage struct {
	ChatID  *int   `json:"chat_id"`
	Message string `json:"message"`
	Silent  bool   `json:"silent"`
}

func NewOutboundTelegramMessage(ChatID *int) OutboundTelegramMessage {
	return OutboundTelegramMessage{
		ChatID:              ChatID,
		DisablePreview:      true,
		DisableNotification: true,
	}
}

func NewInboundTelegramMessage(ChatID *int) InboundTelegramMessage {
	return InboundTelegramMessage{
		ChatID: ChatID,
		Silent: true,
	}
}
