package assistant

import (
	"fmt"
	"github.com/spandigitial/codeassistant/model"
)

func (a ChatGPTBasedCodeAssistant) Article(term string, handlers ...func(markdown string)) error {

	messages := []model.Prompt{
		{
			Role:    "system",
			Content: "Your are a generator of articles about technology, response should have more then one paragraph. Send all responses in Markdown.",
		},
		{
			Role:    "user",
			Content: fmt.Sprintf("Write an article about %s?", term),
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
