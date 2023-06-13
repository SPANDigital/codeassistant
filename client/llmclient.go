// SPDX-License-Identifier: MIT

package client

import (
	"github.com/spandigitial/codeassistant/model"
)

type LanguageModel interface {
	String() string
}

type Message struct {
	Delta string
	Type  string
}

type ModelChan chan LanguageModel
type MessageChan chan Message

type LLMClient interface {
	Models(mdoels ModelChan) error
	Completion(command *model.CommandInstance, messages MessageChan) error
}
