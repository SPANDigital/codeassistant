package assistant

import "github.com/spandigitial/codeassistant/model"

type CodeAssistant interface {
	RailsSchemaToEntities(railsSchema string, entityHandlers []func(code model.SourceCode) model.SourceCode, serviceHandlers []func(code model.SourceCode) model.SourceCode) error
}
