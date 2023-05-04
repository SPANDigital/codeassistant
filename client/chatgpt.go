// SPDX-License-Identifier: MIT

package client

import (
	model2 "github.com/spandigitial/codeassistant/client/model"
	"github.com/spandigitial/codeassistant/model"
)

type ChoiceHandler func(objectType string, choice model2.Choice)

type ChatGPT interface {
	Completion(command *model.CommandInstance, handlers ...ChoiceHandler) error
}
