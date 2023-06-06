package model

type VertexAIConfig struct {
	Model           string   `yaml:"model"`
	Temperature     *float32 `yaml:"temperature"`
	MaxOutputTokens *int     `yaml:"max-output-tokens"`
	TopK            *float32 `yaml:"top-k"`
	TopP            *float32 `yaml:"top-p"`
}
