package vertexai

type prediction struct {
	SafetyAttributes safetyAttributes `json:"safetyAttributes"`
	Content          string           `json:"content"`
	CitationMetadata citationMetadata `json:"citationMetadata"`
}
