package vertexai

type safetyAttributes struct {
	Categories []interface{} `json:"categories"`
	Blocked    bool          `json:"blocked"`
	Scores     []interface{} `json:"scores"`
}
