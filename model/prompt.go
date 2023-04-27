// SPDX-License-Identifier: MIT

package model

import (
	"bytes"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
	"strings"
)

type Prompt struct {
	Role    string `json:"role" yaml:"role"`
	Content string `json:"content"" yaml:"content"`
}

func writeLines(w bytes.Buffer, source []byte, n ast.Node) {
	l := n.Lines().Len()
	for i := 0; i < l; i++ {
		line := n.Lines().At(i)
		w.Write(line.Value(source))
	}
}

func (m Prompt) FencedCodeBlocks() []FencedCodeBlock {
	source := []byte(m.Content)
	node := goldmark.DefaultParser().Parse(text.NewReader(source))
	var codeBlocks []FencedCodeBlock
	ast.Walk(node, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
		if node.Kind() == ast.KindFencedCodeBlock {
			var language string
			var content string

			fcb := node.(*ast.FencedCodeBlock)
			if !entering && fcb.Info != nil {
				segment := fcb.Info.Segment
				language = string(source[segment.Start:segment.Stop])
				var sb strings.Builder
				lines := fcb.BaseBlock.Lines()
				l := lines.Len()
				for i := 0; i < l; i++ {
					line := lines.At(i)
					sb.Write(line.Value(source))
				}
				content = sb.String()
				if language != "" && content != "" {
					codeBlocks = append(codeBlocks, FencedCodeBlock{
						Language: language,
						Content:  content,
					})
				}
			}

		}

		return ast.WalkContinue, nil
	})
	return codeBlocks
}
