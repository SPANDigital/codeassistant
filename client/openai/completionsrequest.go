// SPDX-License-Identifier: MIT

package openai

import (
	"github.com/spandigitial/codeassistant/model/prompts"
)

type completionsRequest struct {
	Model       string           `json:"model"`
	Messages    []prompts.Prompt `json:"messages"`
	User        *string          `json:"user"`
	Temperature *float32         `json:"temperature,omitempty"`
	TopP        *float32         `json:"top_p,omitempty"`
	Stream      bool             `json:"stream"`
}
