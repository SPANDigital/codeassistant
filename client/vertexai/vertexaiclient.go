package vertexai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/spandigitial/codeassistant/client"
	"github.com/spandigitial/codeassistant/client/debugger"
	"github.com/spandigitial/codeassistant/model/prompts"
	"github.com/spf13/viper"
	"io"
	"net/http"
	"time"
)

type Client struct {
	accessToken string
	projectId   string
	location    string
	model       string
	debugger    *debugger.Debugger
	httpClient  *http.Client
}

type Option func(client *Client)

func New(projectId string, location string, model string, debugger *debugger.Debugger, options ...Option) *Client {
	accessToken, _ := generateAccessToken()
	c := &Client{
		accessToken: accessToken,
		projectId:   projectId,
		location:    location,
		model:       model,
		debugger:    debugger,
	}

	for _, option := range options {
		option(c)
	}

	if c.httpClient == nil {
		c.httpClient = http.DefaultClient
	}

	return c
}

func (c *Client) Models(models chan<- client.LanguageModel) error {
	return nil
}

func (c *Client) Completion(commandInstance *prompts.CommandInstance, messageParts chan<- client.MessagePart) error {

	temperature := float64(0.2)
	if commandInstance.Command.VertexAIConfig.Temperature != nil {
		temperature = *commandInstance.Command.VertexAIConfig.Temperature
	}
	maxOutputTokens := 256
	if commandInstance.Command.VertexAIConfig.MaxOutputTokens != nil {
		maxOutputTokens = *commandInstance.Command.VertexAIConfig.MaxOutputTokens
	}
	topP := float64(0.8)
	if commandInstance.Command.VertexAIConfig.TopP != nil {
		topP = *commandInstance.Command.VertexAIConfig.TopP
	}
	topK := 40
	if commandInstance.Command.VertexAIConfig.TopK != nil {
		topK = *commandInstance.Command.VertexAIConfig.TopK
	}

	parameters := parameters{
		Temperature:     temperature,
		MaxOutputTokens: maxOutputTokens,
		TopP:            topP,
		TopK:            topK,
	}

	prompt := commandInstance.JoinedPromptsContent("\n")
	request := predictRequest{
		Instances: []map[string]interface{}{{
			viper.GetString("vertexAiPromptAttribute"): prompt,
		}},
		Parameters: parameters,
	}

	c.debugger.Message(debugger.SentPrompt, prompt)

	requestBytes, err := json.Marshal(request)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("https://%s-aiplatform.googleapis.com/v1/projects/%s/locations/%s/publishers/google/models/%s:predict",
		c.location,
		c.projectId,
		c.location,
		c.model)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(requestBytes))
	if err != nil {
		return err
	}

	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.accessToken))

	requestTime := time.Now()
	c.debugger.Message(debugger.RequestTime, fmt.Sprintf("%v", requestTime))
	// Send the HTTP request]
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	firstResponseTime := time.Now()
	c.debugger.Message(debugger.FirstResponseTime, fmt.Sprintf("%v elapsed %v", firstResponseTime, firstResponseTime.Sub(requestTime)))
	c.debugger.Message(debugger.LastResponseTime, fmt.Sprintf("%v elapsed %v", firstResponseTime, firstResponseTime.Sub(requestTime)))

	// Read the response body
	responseBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var response predictResponse
	err = json.Unmarshal(responseBytes, &response)
	if err != nil {
		return err
	}

	messageParts <- client.MessagePart{Delta: "", Type: "Start"}
	for _, prediction := range response.Predictions {
		messageParts <- client.MessagePart{Delta: prediction.Content, Type: "Part"}
	}
	messageParts <- client.MessagePart{Delta: "", Type: "Done"}
	close(messageParts)
	return nil
}

func (c *Client) Embeddings(model string, input string) ([]float32, error) {
	return nil, nil
}
