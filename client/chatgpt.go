package client

type ChatGPT interface {
	Completion(prompt string) string
}
