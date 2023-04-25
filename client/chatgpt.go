package client

import (
	model2 "github.com/spandigitial/codeassistant/client/model"
	"github.com/spandigitial/codeassistant/model"
)

type ChatGPT interface {
	Completion(command *model.CommandInstance) ([]model2.Choice, error)
}
