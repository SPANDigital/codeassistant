package assistant

import (
	"github.com/spandigitial/codeassistant/model"
	"io"
)

type CodeAssistant interface {
	WhatIs(term string, handlers ...func(markdown string)) error
	RailsSchemaToEntities(railsSchema string, entityHandlers []func(code model.SourceCode) model.SourceCode, serviceHandlers []func(code model.SourceCode) model.SourceCode) error
	Convert(rubyCode io.Reader, rubyType string, nestJsType string, codeHandlers []func(code model.SourceCode) model.SourceCode) error
}
