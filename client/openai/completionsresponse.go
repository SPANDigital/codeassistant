// SPDX-License-Identifier: MIT

package openai

type CompletionsResponse struct {
	Id      string   `json:"id"`
	Object  string   `json:"object"`
	Created int      `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Error   *Error   `json:"error"`
}
