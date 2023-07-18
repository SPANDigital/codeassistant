// SPDX-License-Identifier: MIT

package client

import (
	"github.com/spandigitial/codeassistant/model/prompts"
	"github.com/spandigitial/codeassistant/vectors"
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
	Embeddings(model string, input string) (vectors.Vector, error)
	SimpleCompletion(model string, roleHint string, input string, messageParts chan<- MessagePart) error
}
