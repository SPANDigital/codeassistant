package vertexai

import (
	aiplatform "cloud.google.com/go/aiplatform/apiv1"
	"cloud.google.com/go/aiplatform/apiv1/aiplatformpb"
	"context"
	"fmt"
	"github.com/spandigitial/codeassistant/client"
	"github.com/spandigitial/codeassistant/client/debugger"
	"github.com/spandigitial/codeassistant/model"
	"github.com/spf13/viper"
	"google.golang.org/api/option"
	"google.golang.org/protobuf/types/known/structpb"
	"time"
)

type Client struct {
	projectId string
	location  string
	debugger  *debugger.Debugger
}

type Option func(client *Client)

func New(projectId string, location string, debugger *debugger.Debugger, options ...Option) *Client {
	c := &Client{
		projectId: projectId,
		location:  location,
		debugger:  debugger,
	}

	for _, option := range options {
		option(c)
	}
	return c
}

func (c *Client) Models(models chan<- client.LanguageModel) error {
	return nil
}

func (c *Client) Completion(commandInstance *model.CommandInstance, messageParts chan<- client.MessagePart) error {
	ctx := context.Background()
	pc, err := aiplatform.NewPredictionClient(ctx,
		option.WithEndpoint(fmt.Sprintf("%s-aiplatform.googleapis.com:443", c.location)))
	if err != nil {
		return err
	}
	tctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel() // Always call cancel.

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
		Endpoint: fmt.Sprintf("projects/%s/locations/%s/endpoints/%s",
			viper.GetString("vertexAiProjectId"),
			viper.GetString("vertexAiLocation"),
			viper.GetString("vertexAiModel"),
		),
		Instances:  instances,
		Parameters: parameters,
	}

	resp, err := pc.Predict(tctx, req)
	if err != nil {
		return err
	}
	messageParts <- client.MessagePart{Delta: "", Type: "Start"}
	for _, prediction := range resp.Predictions {
		messageParts <- client.MessagePart{Delta: prediction.GetStructValue().Fields["Content"].GetStringValue(), Type: "Part"}
	}
	messageParts <- client.MessagePart{Delta: "", Type: "Done"}
	return nil
}
