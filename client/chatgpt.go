package client

import (
	model2 "github.com/spandigitial/codeassistant/client/model"
	"github.com/spandigitial/codeassistant/model"
)

type ChatGPT interface {
	Completion(messages ...model.Prompt) ([]model2.Choice, error)
}
