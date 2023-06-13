package vertexai

import (
	aiplatform "cloud.google.com/go/aiplatform/apiv1"
	"cloud.google.com/go/aiplatform/apiv1/aiplatformpb"
	"context"
	"fmt"
	"github.com/spandigitial/codeassistant/client"
	"github.com/spandigitial/codeassistant/client/debugger"
	"github.com/spandigitial/codeassistant/model"
	"google.golang.org/protobuf/types/known/structpb"
)

type VertexAiClient struct {
	projectId string
	location  string
	debugger  *debugger.Debugger
}

type Option func(client *VertexAiClient)

func New(projectId string, location string, debugger *debugger.Debugger, options ...Option) *VertexAiClient {
	c := &VertexAiClient{
		projectId: projectId,
		location:  location,
		debugger:  debugger,
	}

	for _, option := range options {
		option(c)
	}
	return c
}

func (c *VertexAiClient) Completion(commandInstance *model.CommandInstance, messages client.MessageChan) error {
	ctx := context.Background()
	pc, err := aiplatform.NewPredictionClient(ctx)
	if err != nil {
		return err
	}
	defer pc.Close()

	temperature := float32(0.2)
	if commandInstance.Command.VertexAIConfig.Temperature != nil {
		temperature = *commandInstance.Command.VertexAIConfig.Temperature
	}
	maxOutputTokens := 256
	if commandInstance.Command.VertexAIConfig.MaxOutputTokens != nil {
		maxOutputTokens = *commandInstance.Command.VertexAIConfig.MaxOutputTokens
	}
	topP := float32(0.8)
	if commandInstance.Command.VertexAIConfig.TopP != nil {
		topP = *commandInstance.Command.VertexAIConfig.TopP
	}
	topK := 40
	if commandInstance.Command.VertexAIConfig.TopK != nil {
		topK = *commandInstance.Command.VertexAIConfig.TopK
	}

	parameters, err := structpb.NewValue(map[string]interface{}{
		"temperature":     temperature,
		"maxOutputTokens": maxOutputTokens,
		"topP":            topP,
		"topK":            topK,
	})
	if err != nil {
		return err
	}
	instances := make([]*structpb.Value, len(commandInstance.Prompts))
	for idx, prompt := range commandInstance.Prompts {
		instances[idx], err = structpb.NewValue(map[string]interface{}{
			"content": prompt.Content,
		})
		if err != nil {
			return err
		}
	}

	req := &aiplatformpb.PredictRequest{
		Endpoint:   fmt.Sprintf("%s-aiplatform.googleapis.com:443", c.location),
		Instances:  instances,
		Parameters: parameters,
	}

	resp, err := pc.Predict(ctx, req)
	if err != nil {
		return err
	}
	messages <- client.Message{Delta: "", Type: "Start"}
	for _, prediction := range resp.Predictions {
		messages <- client.Message{Delta: prediction.GetStructValue().Fields["Content"].GetStringValue(), Type: "Part"}
	}
	messages <- client.Message{Delta: "", Type: "Done"}
	return nil
}
