// SPDX-License-Identifier: MIT

package openai

import (
	"github.com/spandigitial/codeassistant/model/prompts"
)

type choice struct {
	Delta        *delta          `json:"delta"`
	Message      *prompts.Prompt `json:"message"`
	FinishReason string          `json:"finish_reason"`
	Index        int             `json:"index"`
}

func (c choice) String() string {
	return c.Message.Content
}
