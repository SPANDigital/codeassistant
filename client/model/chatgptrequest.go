// SPDX-License-Identifier: MIT

package model

import (
	"github.com/spandigitial/codeassistant/model"
)

type ChatGPTRequest struct {
	Model       string         `json:"model"`
	Messages    []model.Prompt `json:"messages"`
	User        *string        `json:"user"`
	Temperature *float32       `json:"temperature"`
	TopP        *float32       `json:"top_p"`
	Stream      bool           `json:"stream"`
}
