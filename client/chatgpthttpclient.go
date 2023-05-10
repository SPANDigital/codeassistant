// SPDX-License-Identifier: MIT

package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/spandigitial/codeassistant/client/debugger"
	model2 "github.com/spandigitial/codeassistant/client/model"
	"github.com/spandigitial/codeassistant/model"
	"github.com/spandigitial/codeassistant/ratelimit"
	"golang.org/x/time/rate"
	"io"
	"net/http"
	"regexp"
	"time"
)

type ChatGPTHttpClient struct {
	apiKey       string
	debugger     *debugger.Debugger
	rateLimiter  *rate.Limiter
	httpClient   *http.Client
	rlHTTPClient *ratelimit.RLHTTPClient
	user         *string
	userAgent    *string
}

type Option func(client *ChatGPTHttpClient)

func New(apiKey string, debugger *debugger.Debugger, rateLimiter *rate.Limiter, options ...Option) *ChatGPTHttpClient {
	c := &ChatGPTHttpClient{
		apiKey:      apiKey,
		debugger:    debugger,
		rateLimiter: rateLimiter,
	}

	for _, option := range options {
		option(c)
	}

	if c.httpClient == nil {
		c.httpClient = http.DefaultClient
	}

	c.rlHTTPClient = &ratelimit.RLHTTPClient{
		Client:      c.httpClient,
		Ratelimiter: c.rateLimiter,
	}
	return c
}

func WithHttpClient(httpClient *http.Client) Option {
	return func(client *ChatGPTHttpClient) {
		client.httpClient = httpClient
	}
}

func WithUser(user string) Option {
	return func(client *ChatGPTHttpClient) {
		client.user = &user
	}
}

func WithUserAgent(userAgent string) Option {
	return func(client *ChatGPTHttpClient) {
		client.userAgent = &userAgent
	}
}

var dataRegex = regexp.MustCompile("data: (\\{.+\\})\\w?")

func (c *ChatGPTHttpClient) Completion(commandInstance *model.CommandInstance, handlers ...ChoiceHandler) error {
	url := "https://api.openai.com/v1/chat/completions"

	for _, prompt := range commandInstance.Prompts {
		c.debugger.Message("sent-prompt", fmt.Sprintf("(%s) %s", prompt.Role, prompt.Content))
	}

	// Create the request body
	request := model2.ChatGPTRequest{
		Messages: commandInstance.Prompts,
		User:     c.user,
		Stream:   true,
	}

	if commandInstance.Command.Model != "" {
		request.Model = commandInstance.Command.Model
	} else {
		request.Model = "gpt-3.5-turbo"
	}
	if commandInstance.Command.Temperature != nil {
		request.Temperature = commandInstance.Command.Temperature
	}
	if commandInstance.Command.TopP != nil {
		request.TopP = commandInstance.Command.TopP
	}

	requestBytes, err := json.Marshal(request)
	if err != nil {
		panic(err)
	}

	if c.debugger.IsRecording("request-payload") {
		c.debugger.Message("request-payload", string(requestBytes))
	}

	if c.debugger.IsRecording("request-tokens") {
		c.debugger.Message("request-tokens", fmt.Sprintf("%d", debugger.NumTokensFromRequest(request)))
	}

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
	resp, err := c.rlHTTPClient.Do(req)
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
				var response model2.ChatGPTResponse
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
					var response model2.ChatGPTResponse
					err = json.Unmarshal(matches[1], &response)
					if response.Error != nil {
						return response.Error
					}
					if err == nil {
						for _, choice := range response.Choices {
							for _, handler := range handlers {
								handler(response.Object, choice)
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

	return nil
}
