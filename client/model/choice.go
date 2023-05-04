// SPDX-License-Identifier: MIT

package model

import "github.com/spandigitial/codeassistant/model"

type Choice struct {
	Delta        *Delta        `json:"delta"`
	Message      *model.Prompt `json:"message"`
	FinishReason string        `json:"finish_reason"`
	Index        int           `json:"index"`
}

func (c Choice) String() string {
	return c.Message.Content
}
