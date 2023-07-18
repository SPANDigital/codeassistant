package indexing

import (
	"path/filepath"
	"strings"
)

type Library struct {
	Name           string
	DisplayName    string `yaml:"display-name"`
	Icon           string
	Path           string   `json:"-"`
	BuiltFromPaths []string `json:"-"`
	Index          string
	Data           map[string]interface{} `yaml:"-" json:"-"'`
	Embeddings     map[string]*Indexing
}

func (l *Library) addBuildPath(path string) {
	l.BuiltFromPaths = append(l.BuiltFromPaths, filepath.Base(path))
}

func (l *Library) getEmbedding(path string) *Indexing {
	base := filepath.Base(path)
	frontName := strings.Split(base, ".")[0]
	command, found := l.Embeddings[frontName]
	if !found {
		command = &Indexing{
			Name:    frontName,
			Library: l,
		}
		l.Embeddings[frontName] = command
	}
	command.BuiltFromPaths = append(command.BuiltFromPaths, base)
	return command
}
