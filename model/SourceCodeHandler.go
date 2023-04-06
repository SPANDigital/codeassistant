package model

type SourceCodeHandler func(code SourceCode) SourceCode

func SourceCodeHandlers(handlers ...SourceCodeHandler) []SourceCodeHandler {
	return handlers
}

func NoSourceCodeHandlers() []SourceCodeHandler {
	return []SourceCodeHandler{}
}
