package model

type OpenAIConfig struct {
	Model       string   `yaml:"model"`
	Temperature *float32 `yaml:"temperature"`
	TopP        *float32 `yaml:"top_p"`
}
