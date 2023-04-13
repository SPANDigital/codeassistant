package model

type FencedCodeBlock struct {
	Language string
	Content  string
}

func (b FencedCodeBlock) ToSourceCode(filenameGenerator func(block FencedCodeBlock) string) SourceCode {
	return SourceCode{
		Filename: filenameGenerator(b),
		Language: b.Language,
		Content:  b.Content,
	}
}

func (b FencedCodeBlock) String() string {
	return b.Content
}
