package model

import "github.com/spandigitial/codeassistant/model"

type Choice struct {
	Message      model.Message `json:"message"`
	FinishReason string        `json:"finish_reason"`
	Index        int           `json:"index"`
}
