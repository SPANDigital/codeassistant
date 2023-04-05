package client

import (
	"bytes"
	"encoding/json"
	"github.com/jpfielding/gowirelog/wirelog"
	"github.com/spandigitial/codeassistant/model"
	"io"
	"net/http"
	"os"
)

type ChatGPTHttpClient struct {
	client *http.Client
	apiKey string
}

func New(apiKey string) *ChatGPTHttpClient {
	transport := wirelog.NewHTTPTransport()
	wirelog.LogToWriter(transport, os.Stdout, true, true)
	return &ChatGPTHttpClient{
		apiKey: apiKey,
		client: &http.Client{
			Transport: transport,
		},
	}
}

func (c *ChatGPTHttpClient) Completion(prompt string) string {
	url := "https://api.openai.com/v1/chat/completions"

	// Create the request body
	request := model.ChatGPTRequest{
		Prompt: prompt,
		Model:  "gpt-3.5-turbo",
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
		panic(err)
	}
	defer resp.Body.Close()

	// Read the response body
	responseBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	// Parse the response JSON
	var response model.ChatGPTResponse
	err = json.Unmarshal(responseBytes, &response)
	if err != nil {
		panic(err)
	}

	return response.Completion
}
