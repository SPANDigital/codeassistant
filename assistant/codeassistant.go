package assistant

import "github.com/spandigitial/codeassistant/model"

type CodeAssistant interface {
	RailsSchemaToEntities(railsSchema string, handlers ...func(code model.SourceCode) model.SourceCode) error
}
