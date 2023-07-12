package prompts

type UiHint struct {
	Label      string                 `yaml:"label"`
	HelperText string                 `yaml:"helper-text"`
	Props      map[string]interface{} `yaml:"props"`
}
