package domain

// Message decoded from json
type UserMessage struct {
	ClientMessageID string `json:"clientMessageId"`
	RecipientID     string `json:"recipientId,omitempty"`
	ChatID          string `json:"chatId,omitempty"`
	Content         string `json:"content"`
}

// Message that client sends to hub
type TargetedMessage struct {
	Content      []byte
	RecipientIDs []string
}
