package assistant

import (
	"github.com/spandigitial/codeassistant/client"
)

type Option func()

type ChatGPTBasedCodeAssistant struct {
	chatGPT client.ChatGPT
}

func New(chatGPT client.ChatGPT) *ChatGPTBasedCodeAssistant {
	return &ChatGPTBasedCodeAssistant{
		chatGPT: chatGPT,
	}
}
