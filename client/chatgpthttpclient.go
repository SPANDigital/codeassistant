// SPDX-License-Identifier: MIT

package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	model2 "github.com/spandigitial/codeassistant/client/model"
	"github.com/spandigitial/codeassistant/model"
	"github.com/spandigitial/codeassistant/ratelimit"
	"golang.org/x/time/rate"
	"io"
	"net/http"
	"os"
	"regexp"
)

type ChatGPTHttpClient struct {
	apiKey       string
	rateLimiter  *rate.Limiter
	httpClient   *http.Client
	rlHTTPClient *ratelimit.RLHTTPClient
	user         *string
}

type Option func(client *ChatGPTHttpClient)

func New(apiKey string, rateLimiter *rate.Limiter, options ...Option) *ChatGPTHttpClient {
	c := &ChatGPTHttpClient{
		apiKey:      apiKey,
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

var dataRegex = regexp.MustCompile("data: (.+)")

func (c *ChatGPTHttpClient) Completion(commandInstance *model.CommandInstance, handlers ...ChoiceHandler) error {
	url := "https://api.openai.com/v1/chat/completions"

	fmt.Fprintln(os.Stderr, "Send following prompts to ChatGPT")
	for _, prompt := range commandInstance.Prompts {
		fmt.Fprintf(os.Stderr, ">>> (%s) %s\n", prompt.Role, prompt.Content)
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

	// Create the HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBytes))
	if err != nil {
		panic(err)
	}
	req.Header.Set("User-Agent", "SPAN Digital code assistant")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	// Send the HTTP request]
	resp, err := c.rlHTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	for {
		data := make([]byte, 1024)
		read, err := resp.Body.Read(data)
		fmt.Fprintln(os.Stderr, "Response", string(data))
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		if len(data) > 0 && string(data[:1]) == "{" {
			var response model2.ChatGPTResponse
			err = json.Unmarshal(data[:read], &response)
			if response.Error != nil {
				return response.Error
			}
			if err != nil {
				return err
			}
		}

		matches := dataRegex.FindSubmatch(data)

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

	return nil
}
