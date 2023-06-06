package vertexai

import "github.com/spandigitial/codeassistant/client/debugger"

type VertexAiClient struct {
	apiKey   string
	location string
	debugger *debugger.Debugger
}

type Option func(client *VertexAiClient)

func New(apiKey string, location string, debugger *debugger.Debugger, options ...Option) *VertexAiClient {
	c := &VertexAiClient{
		apiKey:   apiKey,
		location: location,
		debugger: debugger.
	}

	for _, option := range options {
		option(c)
	}
}
