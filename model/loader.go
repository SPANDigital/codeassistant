package model

type Loader interface {
	Load(records <-chan map[string]interface{}) error
}
