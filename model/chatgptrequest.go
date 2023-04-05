package model

type ChatGPTRequest struct {
	Prompt string `json:"prompt"`
	Model  string `json:"model"`
}
