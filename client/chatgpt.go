package client

type ChatGPT interface {
	completion(prompt string) string
}
