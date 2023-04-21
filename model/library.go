package model

type Library struct {
	Name     string
	FullPath string
	Index    string
	Commands map[string]*Command
}
