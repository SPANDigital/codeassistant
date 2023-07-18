package indexing

import "github.com/spandigitial/codeassistant/vectors"

type IndexItem struct {
	Content map[string]interface{}
	Vectors map[string]vectors.Vector
}
