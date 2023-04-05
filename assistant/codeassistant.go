package assistant

type CodeAssistant interface {
	RailsSchemaToEntities(railsSchema string) []Code
}
