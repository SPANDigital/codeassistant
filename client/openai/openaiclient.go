// SPDX-License-Identifier: MIT

package openai

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/spandigitial/codeassistant/client"
	"github.com/spandigitial/codeassistant/client/debugger"
	"github.com/spandigitial/codeassistant/model/prompts"
	"github.com/spf13/viper"
	"golang.org/x/time/rate"
	"io"
	"net/http"
	"regexp"
	"time"
)

type OpenAiClient struct {
	apiKey      string
	debugger    *debugger.Debugger
	rateLimiter *rate.Limiter
	httpClient  *http.Client
	user        *string
	userAgent   *string
}

type Option func(client *OpenAiClient)

func New(apiKey string, debugger *debugger.Debugger, options ...Option) *OpenAiClient {
	c := &OpenAiClient{
		apiKey:   apiKey,
		debugger: debugger,
	}

	for _, option := range options {
		option(c)
	}

	if c.httpClient == nil {
		c.httpClient = http.DefaultClient
	}

	return c
}

func WithHttpClient(httpClient *http.Client) Option {
	return func(client *OpenAiClient) {
		client.httpClient = httpClient
	}
}

func WithUser(user string) Option {
	return func(client *OpenAiClient) {
		client.user = &user
	}
}

func WithUserAgent(userAgent string) Option {
	return func(client *OpenAiClient) {
		client.userAgent = &userAgent
	}
}

var dataRegex = regexp.MustCompile("data: (\\{.+\\})\\w?")

func (c *OpenAiClient) Models(models chan<- client.LanguageModel) error {
	url := fmt.Sprintf("%s/v1/models", viper.GetString("openAiUrlPrefix"))
	requestTime := time.Now()

	c.debugger.MessageF(debugger.RequestTime, "%v", requestTime)

	// Create the HTTP request
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", *c.userAgent)
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	responseTime := time.Now()
	elapsed := responseTime.Sub(requestTime)
	c.debugger.MessageF(debugger.FirstResponseTime, "%v elapsed %v", responseTime, elapsed)
	c.debugger.MessageF(debugger.LastResponseTime, "%v elapsed %v", responseTime, elapsed)

	// Read the response body
	responseBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// Parse the response JSON
	var response languageModelsResponse
	err = json.Unmarshal(responseBytes, &response)
	if err != nil {
		return err
	}

	for _, languageModel := range response.Data {
		models <- languageModel
	}

	close(models)

	return nil
}

func (c *OpenAiClient) Completion(commandInstance *prompts.CommandInstance, messageParts chan<- client.MessagePart) error {
	url := fmt.Sprintf("%s/v1/chat/prompts", viper.GetString("openAiUrlPrefix"))

	for _, prompt := range commandInstance.Prompts {
		c.debugger.Message(debugger.SentPrompt, fmt.Sprintf("(%s) %s", prompt.Role, prompt.Content))
	}

	// Create the request body
	request := completionsRequest{
		Messages: commandInstance.Prompts,
		User:     c.user,
		Stream:   true,
	}

	if commandInstance.Command.OpenAIConfig.Model != "" {
		request.Model = commandInstance.Command.OpenAIConfig.Model
	} else {
		model := viper.GetString("openAiModel")
		if model == "" {
			model = "gpt-4"
		}
		request.Model = model
	}
	if commandInstance.Command.OpenAIConfig.Temperature != nil {
		request.Temperature = commandInstance.Command.OpenAIConfig.Temperature
	}
	if commandInstance.Command.OpenAIConfig.TopP != nil {
		request.TopP = commandInstance.Command.OpenAIConfig.TopP
	}

	requestBytes, err := json.Marshal(request)
	if err != nil {
		panic(err)
	}

	c.debugger.MessageBytes(debugger.RequestPayload, requestBytes)

	requestTime := time.Now()

	c.debugger.Message(debugger.RequestHeader, fmt.Sprintf("%v", requestTime))

	// Create the HTTP request
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(requestBytes))
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", *c.userAgent)
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	// Send the HTTP request]
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	first := true

	var buff bytes.Buffer
	messageParts <- client.MessagePart{Delta: "", Type: "Start"}
	for {
		data := make([]byte, 1024)
		read, err := resp.Body.Read(data)

		if first {
			firstResponseTime := time.Now()
			c.debugger.Message(debugger.FirstResponseTime, fmt.Sprintf("%v elapsed %v", firstResponseTime, firstResponseTime.Sub(requestTime)))
			first = false
		}

		if err == io.EOF {
			lastResponseTime := time.Now()
			c.debugger.Message(debugger.LastResponseTime, fmt.Sprintf("%v elapsed %v", lastResponseTime, lastResponseTime.Sub(requestTime)))
			return nil
		}
		if err != nil {
			return err
		}

		buff.Write(data[:read])
		bytes := buff.Bytes()
		size := len(bytes)
		if string(bytes[size-1:size]) == "\n" {
			if len(data) > 0 && string(data[:1]) == "{" {
				var response completionsResponse
				err = json.Unmarshal(data[:read], &response)
				if response.Error != nil {
					return response.Error
				}
				if err != nil {
					return err
				}
				return fmt.Errorf("unexecpted response: %s", string(bytes))
			}

			allMatches := dataRegex.FindAllSubmatch(bytes, -1)
			for _, matches := range allMatches {

				if len(matches) > 0 {
					var response completionsResponse
					err = json.Unmarshal(matches[1], &response)
					if response.Error != nil {
						return response.Error
					}
					if err == nil {
						for _, choice := range response.Choices {

							if response.Object == "chat.completion.chunk" && choice.Delta != nil {
								messageParts <- client.MessagePart{Delta: choice.Delta.Content, Type: "Part"}
							}
						}
					} else {
						return err
					}
				}
			}
			buff.Reset()
		}
	}
	messageParts <- client.MessagePart{Delta: "", Type: "Done"}
	close(messageParts)

	return nil
}

func (c *OpenAiClient) Embeddings(model string, input string) ([]float32, error) {
	url := fmt.Sprintf("%s/v1/embeddings", viper.GetString("openAiUrlPrefix"))
	c.debugger.Message(debugger.SentInput, input)
	request := map[string]string{
		"model": model,
		"input": input,
	}

	requestBytes, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	c.debugger.MessageBytes(debugger.RequestPayload, requestBytes)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(requestBytes))
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", *c.userAgent)
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	requestTime := time.Now()
	c.debugger.Message(debugger.RequestTime, fmt.Sprintf("%v", requestTime))
	// Send the HTTP request]
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	firstResponseTime := time.Now()
	c.debugger.Message(debugger.FirstResponseTime, fmt.Sprintf("%v elapsed %v", firstResponseTime, firstResponseTime.Sub(requestTime)))
	c.debugger.Message(debugger.LastResponseTime, fmt.Sprintf("%v elapsed %v", firstResponseTime, firstResponseTime.Sub(requestTime)))

	// Read the response body
	responseBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response embbedingResponse
	err = json.Unmarshal(responseBytes, &response)
	if err != nil {
		return nil, err
	}

	if response.Error != nil {
		return nil, errors.New(response.Error.Message)
	}

	return response.Data[0].Embedding, nil

}
