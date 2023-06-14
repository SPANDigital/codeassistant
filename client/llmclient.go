// SPDX-License-Identifier: MIT

package client

import (
	"github.com/spandigitial/codeassistant/model"
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
	Completion(command *model.CommandInstance, messageParts chan<- MessagePart) error
}
