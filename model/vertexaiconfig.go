package model

type VertexAIConfig struct {
	Model           string   `yaml:"model"`
	Temperature     *float64 `yaml:"temperature"`
	MaxOutputTokens *int     `yaml:"max-output-tokens"`
	TopK            *int     `yaml:"top-k"`
	TopP            *float64 `yaml:"top-p"`
}
