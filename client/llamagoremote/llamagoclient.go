package llamagoremote

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/spandigitial/codeassistant/client"
	"github.com/spandigitial/codeassistant/client/debugger"
	"github.com/spandigitial/codeassistant/model"
	"github.com/spandigitial/codeassistant/transport"
	"io"
	"net/http"
	"time"
)

type Client struct {
	protocol     Protocol
	host         string
	port         int
	debugger     *debugger.Debugger
	httpClient   *http.Client
	userAgent    *string
	pollDuration time.Duration
}

type Option func(client *Client)

func New(protocol Protocol, host string, port int, pollDuration time.Duration, debugger *debugger.Debugger, options ...Option) *Client {
	c := &Client{
		protocol:     protocol,
		host:         host,
		port:         port,
		debugger:     debugger,
		pollDuration: pollDuration,
	}

	for _, option := range options {
		option(c)
	}

	if c.httpClient == nil {
		c.httpClient = http.DefaultClient
	}

	c.httpClient.Transport = transport.New(c.httpClient.Transport, c.debugger)

	return c
}

func WithHttpClient(httpClient *http.Client) Option {
	return func(client *Client) {
		client.httpClient = httpClient
	}
}

func (c *Client) Models(models chan<- client.LanguageModel) error {
	close(models)
	return nil
}

func (c *Client) Completion(ci *model.CommandInstance, messageParts chan<- client.MessagePart) error {

	sendURL := fmt.Sprintf("%s://%s:%d/jobs", c.protocol, c.host, c.port)
	uuid := uuid.New()
	requestTime := time.Now()

	c.debugger.Message("request-time", fmt.Sprintf("%v", requestTime))

	request := promptRequest{
		ID:     uuid,
		Prompt: ci.JoinedPromptsContent("\n\n"),
	}

	c.debugger.Message("sent-prompt", request.Prompt)

	requestBytes, err := json.Marshal(request)
	if err != nil {
		return err
	}

	// Create the HTTP request
	req, err := http.NewRequest("POST", sendURL, bytes.NewBuffer(requestBytes))
	if err != nil {
		return err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	// Read the response body
	responseBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var promptResponse promptResponse
	err = json.Unmarshal(responseBytes, &promptResponse)
	if err != nil {
		return err
	}

	if promptResponse.Status == processing {
		for {
			statusURL := fmt.Sprintf("%s://%s:%d/jobs/status/%s", c.protocol, c.host, c.port, uuid.String())

			req, err := http.NewRequest("GET", statusURL, nil)
			if err != nil {
				return err
			}
			resp, err := c.httpClient.Do(req)
			if err != nil {
				return err
			}
			responseBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			var jobStatusResponse jobStatusResponse
			err = json.Unmarshal(responseBytes, &jobStatusResponse)
			if err != nil {
				return err
			}

			if jobStatusResponse.Status == finished {
				break
			}
			time.Sleep(c.pollDuration)
		}

	}

	jobUrl := fmt.Sprintf("%s://%s:%d/jobs/%s", c.protocol, c.host, c.port, uuid.String())
	req, err = http.NewRequest("GET", jobUrl, nil)

	resp, err = c.httpClient.Do(req)
	if err != nil {
		return err
	}

	responseBytes, err = io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var jobResponse jobResponse
	err = json.Unmarshal(responseBytes, &jobResponse)
	if err != nil {
		return err
	}

	messageParts <- client.MessagePart{Delta: "", Type: "Start"}
	messageParts <- client.MessagePart{Delta: jobResponse.Output, Type: "Part"}
	messageParts <- client.MessagePart{Delta: "", Type: "Done"}
	close(messageParts)

	return nil

}
