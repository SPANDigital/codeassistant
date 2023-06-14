package vertexai

type predictRequest struct {
	Instances  []instance `json:"instances"`
	Parameters parameters `json:"parameters"`
}
