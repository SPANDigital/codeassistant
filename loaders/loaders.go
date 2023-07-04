package loaders

import (
	"github.com/spandigitial/codeassistant/model"
	"path"
)

func New(filename string) model.Loader {
	switch path.Ext(filename) {
	case "csv":
		return &CsvLoader{
			filename: filename,
		}

	}
}
