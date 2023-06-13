// SPDX-License-Identifier: MIT

package client

import (
	model2 "github.com/spandigitial/codeassistant/client/model"
	"github.com/spandigitial/codeassistant/model"
)

type Message struct {
	Delta string
	Type  string
}

type MessageChan chan Message

type ModelHandler func(languageModel model2.LanguageModel)

type LLMClient interface {
	Models(handlers ...ModelHandler) error
	Completion(command *model.CommandInstance, messages MessageChan) error
}
