package debugger

import (
	"fmt"
	"github.com/pkoukk/tiktoken-go"
	"github.com/spandigitial/codeassistant/client/model"
	"strings"
)

// Based upon: https://github.com/openai/openai-cookbook/blob/main/examples/How_to_count_tokens_with_tiktoken.ipynb

func NumTokensFromRequest(request model.ChatGPTRequest) (numTokens int) {
	tkm, err := tiktoken.EncodingForModel(request.Model)
	if err != nil {
		err = fmt.Errorf("EncodingForModel: %v", err)
		fmt.Println(err)
		return
	}

	var tokensPerMessage = 3 // Default to 3

	if strings.HasPrefix(request.Model, "gpt-3.5") {
		tokensPerMessage = 4
	} else if strings.HasPrefix(request.Model, "gpt-3.4") {
		tokensPerMessage = 3
	}

	for _, message := range request.Messages {
		numTokens += tokensPerMessage
		numTokens += len(tkm.Encode(message.Content, nil, nil))
		numTokens += len(tkm.Encode(message.Role, nil, nil))
		// TODO: count tokens for name/role, requires changes in ChatGPTRequest and ChatGPTResponse
	}

	numTokens += 3 // Every reply is primed with <|start|>assistant<|message|>
	return numTokens
}
