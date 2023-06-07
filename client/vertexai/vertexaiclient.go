package vertexai

import (
	aiplatform "cloud.google.com/go/aiplatform/apiv1"
	"cloud.google.com/go/aiplatform/apiv1/aiplatformb"
	"cloud.google.com/go/aiplatform/apiv1/aiplatformpb"
	"context"
	"fmt"
	"github.com/spandigitial/codeassistant/client"
	"github.com/spandigitial/codeassistant/client/debugger"
	"github.com/spandigitial/codeassistant/model"
	"google.golang.org/protobuf/types/known/structpb"
)

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
	return c
}

func (c *VertexAiClient) Completion(commandInstance *model.CommandInstance, handlers ...client.ChoiceHandler) error {
	ctx := context.Background()
	pc, err := aiplatform.NewPredictionClient(ctx)
	if err != nil {
		return err
	}
	defer pc.Close()

	parameters, err := structpb.NewValue(map[string]interface{}{
	}
	if err != nil {
		return err
	}
	instances := make([]*structpb.Value, len(commandInstance.Prompts))
    for idx, prompt := range commandInstance.Prompts {
		instances[idx], err = structpb.NewValue(map[string]interface{}{
			"content": prompt.Content,
		}
		if err != nil {
			return err
		}
	}


	req := &aiplatformpb.PredictRequest{
		Endpoint:   fmt.Sprintf("%s-aiplatform.googleapis.com:443", c.location),
		Instances:  instances,
		Parameters: parameters,
	}


	op, err := aic.CreateDataset(ctx, req)
	if err != nil {
		return err
	}

	resp, err := op.Wait(ctx)
	if err != nil {
		return err
	}

}
