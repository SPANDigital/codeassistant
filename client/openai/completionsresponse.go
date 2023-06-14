// SPDX-License-Identifier: MIT

package openai

type completionsResponse struct {
	Id      string       `json:"id"`
	Object  string       `json:"object"`
	Created int          `json:"created"`
	Model   string       `json:"model"`
	Choices []choice     `json:"choices"`
	Error   *openAiError `json:"error"`
}
