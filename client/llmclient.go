// SPDX-License-Identifier: MIT

package client

import (
	"github.com/spandigitial/codeassistant/model/prompts"
)

type LanguageModel interface {
	String() string
}

type MessagePart struct {
	Delta string
	Type  string
}

type LLMClient interface {
	Models(models chan<- LanguageModel) error
	Completion(command *prompts.CommandInstance, messageParts chan<- MessagePart) error
	Embeddings(model string, input string) ([]float32, error)
}
