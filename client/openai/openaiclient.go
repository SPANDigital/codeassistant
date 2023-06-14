// SPDX-License-Identifier: MIT

package openai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/spandigitial/codeassistant/client"
	"github.com/spandigitial/codeassistant/client/debugger"
	"github.com/spandigitial/codeassistant/model"
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
	url := "https://api.openai.com/v1/models"
	requestTime := time.Now()

	c.debugger.Message("request-time", fmt.Sprintf("%v", requestTime))

	// Create the HTTP request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", *c.userAgent)
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	if c.debugger.IsRecording("request-header") {
		var bytes bytes.Buffer
		req.Header.Write(&bytes)
		c.debugger.Message("request-header", bytes.String())
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if c.debugger.IsRecording("response-header") {
		var bytes bytes.Buffer
		resp.Header.Write(&bytes)
		c.debugger.Message("response-header", bytes.String())
	}

	responseTime := time.Now()
	elapsed := responseTime.Sub(requestTime)
	c.debugger.Message("first-response-time", fmt.Sprintf("%v elapsed %v", responseTime, elapsed))
	c.debugger.Message("last-response-time", fmt.Sprintf("%v elapsed %v", responseTime, elapsed))

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
	return nil
}

func (c *OpenAiClient) Completion(commandInstance *model.CommandInstance, messageParts chan<- client.MessagePart) error {
	url := "https://api.openai.com/v1/chat/completions"

	for _, prompt := range commandInstance.Prompts {
		c.debugger.Message("sent-prompt", fmt.Sprintf("(%s) %s", prompt.Role, prompt.Content))
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
		model := viper.GetString("defaultOpenAiModel")
		if model == "" {
			model = "gpt-3.5-turbo"
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

	if c.debugger.IsRecording("request-payload") {
		c.debugger.Message("request-payload", string(requestBytes))
	}

	/*
		if c.debugger.IsRecording("request-tokens") {
			c.debugger.MessagePart("request-tokens", fmt.Sprintf("%d", debugger.NumTokensFromRequest(request)))
		}
	*/

	requestTime := time.Now()

	c.debugger.Message("request-time", fmt.Sprintf("%v", requestTime))

	// Create the HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBytes))
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", *c.userAgent)
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	if c.debugger.IsRecording("request-header") {
		var bytes bytes.Buffer
		req.Header.Write(&bytes)
		c.debugger.Message("request-header", bytes.String())
	}

	// Send the HTTP request]
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	if c.debugger.IsRecording("response-header") {
		var bytes bytes.Buffer
		resp.Header.Write(&bytes)
		c.debugger.Message("response-header", bytes.String())
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
			c.debugger.Message("first-response-time", fmt.Sprintf("%v elapsed %v", firstResponseTime, firstResponseTime.Sub(requestTime)))
			first = false
		}

		if err == io.EOF {
			lastResponseTime := time.Now()
			c.debugger.Message("last-response-time", fmt.Sprintf("%v elapsed %v", lastResponseTime, lastResponseTime.Sub(requestTime)))
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
