package client

import (
	"bytes"
	"encoding/json"
	"github.com/jpfielding/gowirelog/wirelog"
	model2 "github.com/spandigitial/codeassistant/client/model"
	"github.com/spandigitial/codeassistant/model"
	"github.com/spandigitial/codeassistant/ratelimit"
	"golang.org/x/time/rate"
	"io"
	"net/http"
)

type ChatGPTHttpClient struct {
	client *ratelimit.RLHTTPClient
	apiKey string
}

func New(apiKey string, ratelimiter *rate.Limiter) *ChatGPTHttpClient {
	transport := wirelog.NewHTTPTransport()
	wirelog.LogToFile(transport, "/tmp/http.log", true, true)
	return &ChatGPTHttpClient{
		apiKey: apiKey,
		client: &ratelimit.RLHTTPClient{
			Client: &http.Client{
				Transport: transport,
			},
			Ratelimiter: ratelimiter,
		},
	}
}

func (c *ChatGPTHttpClient) Completion(messages ...model.Message) ([]model2.Choice, error) {
	url := "https://api.openai.com/v1/chat/completions"

	// Create the request body
	request := model2.ChatGPTRequest{
		Messages: messages,
		Model:    "gpt-3.5-turbo-0301",
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
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	// Send the HTTP request]
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response body
	responseBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Parse the response JSON
	var response model2.ChatGPTResponse
	err = json.Unmarshal(responseBytes, &response)
	if err != nil {
		return nil, err
	}

	return response.Choices, nil
}
