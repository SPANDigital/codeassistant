package assistant

import (
	"fmt"
	"github.com/spandigitial/codeassistant/model"
)

func (a ChatGPTBasedCodeAssistant) WhatIs(term string, handlers ...func(markdown string)) error {

	messages := []model.Message{
		{
			Role:    "system",
			Content: "Send all responses in Markdown. Favour longer responses.",
		},
		{
			Role:    "user",
			Content: fmt.Sprintf("What is %s?", term),
		},
	}

	choices, err := a.chatGPT.Completion(messages...)
	if err != nil {
		return err
	}

	for _, choice := range choices {
		for _, handler := range handlers {
			handler(choice.Message.Content)
		}
	}
	return nil
}
