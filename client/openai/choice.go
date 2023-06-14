// SPDX-License-Identifier: MIT

package openai

import "github.com/spandigitial/codeassistant/model"

type choice struct {
	Delta        *delta        `json:"delta"`
	Message      *model.Prompt `json:"message"`
	FinishReason string        `json:"finish_reason"`
	Index        int           `json:"index"`
}

func (c choice) String() string {
	return c.Message.Content
}
