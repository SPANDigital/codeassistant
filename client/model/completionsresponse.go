// SPDX-License-Identifier: MIT

package model

type CompletionsResponse struct {
	Id      string   `json:"id"`
	Object  string   `json:"object"`
	Created int      `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Error   *Error   `json:"error"`
}
