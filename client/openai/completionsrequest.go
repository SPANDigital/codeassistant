// SPDX-License-Identifier: MIT

package openai

import (
	"github.com/spandigitial/codeassistant/model"
)

type completionsRequest struct {
	Model       string         `json:"model"`
	Messages    []model.Prompt `json:"messages"`
	User        *string        `json:"user"`
	Temperature *float32       `json:"temperature,omitempty"`
	TopP        *float32       `json:"top_p,omitempty"`
	Stream      bool           `json:"stream"`
}
