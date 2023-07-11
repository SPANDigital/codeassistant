package vertexai

type predictRequest struct {
	Instances  []map[string]interface{} `json:"instances"`
	Parameters parameters               `json:"parameters"`
}
