package indexing

import "github.com/spandigitial/codeassistant/model/ui"

type Indexing struct {
	Library            *Library             `json:"-"`
	Name               string               `yaml:"-"`
	DisplayName        string               `yaml:"display-name"`
	BuiltFromPaths     []string             `yaml:"-" json:"-"`
	Params             map[string]string    `yaml:"params"`
	UiHints            map[string]ui.UiHint `yaml:"ui-hints"`
	Script             string               `yaml:"-"`
	Clone              string               `yaml:"clone"`
	AccessToken        string               `yaml:"access-token"`
	Collect            string               `yaml:"collect"`
	Vectors            map[string]string    `yaml:"vectors"`
	IndexJson          string               `yaml:"index-json"`
	LoadIfPossible     bool                 `yaml:"load-if-possible"`
	SaveIfPossible     bool                 `yaml:"save-if-possible"`
	WorkingDirectory   string               `yaml:"working-directory"`
	ExtractFilesEnrich map[string]string    `yaml:"extract-files-enrich"`
	RoleHint           string               `yaml:"role-hint"`
}
