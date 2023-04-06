package model

type ChatGPTResponse struct {
	Id      string   `json:"id"`
	Object  string   `json:"object"`
	Created int      `json:"Created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
}
