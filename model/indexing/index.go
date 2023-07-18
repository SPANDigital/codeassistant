package indexing

import (
	"github.com/spandigitial/codeassistant/vectors"
	"sort"
)

type Index struct {
	Items []IndexItem
}

func NewIndex(contents []map[string]interface{}) *Index {
	items := make([]IndexItem, len(contents))
	for idx, c := range contents {
		items[idx] = IndexItem{
			Content: c,
			Vectors: make(map[string]vectors.Vector),
		}
	}
	return &Index{Items: items}
}

func (i *Index) Length() int {
	return len(i.Items)
}

func (i *Index) Closest(content string, number int, vector vectors.Vector) ([]IndexItem, error) {
	s := i.Items
	sort.Slice(s, func(i, j int) bool {
		iDistance, err1 := vectors.CosineSimilarity(s[i].Vectors[content], vector)
		jDistance, err2 := vectors.CosineSimilarity(s[j].Vectors[content], vector)
		if err1 != nil || err2 != nil {
			return false
		}
		return iDistance > jDistance
	})
	return s[:number], nil

}
