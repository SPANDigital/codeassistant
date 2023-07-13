package embeddings

import "github.com/spandigitial/codeassistant/model/ui"

type Embedding struct {
	Library        *Library             `json:"-"`
	Name           string               `yaml:"-"`
	DisplayName    string               `yaml:"display-name"`
	BuiltFromPaths []string             `yaml:"-" json:"-"`
	Params         map[string]string    `yaml:"params"`
	UiHints        map[string]ui.UiHint `yaml:"ui-hints"`
	Script         string               `yaml:"-"`
	Clone          string               `yaml:"clone"`
	AccessToken    string               `yaml:"access-token"`
	Build          string               `yaml:"build"`
	Input          string               `yaml:"input"`
}
