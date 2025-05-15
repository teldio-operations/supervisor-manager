package module

import "github.com/invopop/jsonschema"

type Info struct {
	Name        string `json:"name"`
	Title       string `json:"title"`
	Description string `json:"description"`

	Config *jsonschema.Schema `json:"config"`

	Listens []*EventInfo `json:"listens"`
	Emits   []*EventInfo `json:"emits"`
}

type EventInfo struct {
	jsonschema.Schema
	Reply jsonschema.Schema `json:"reply"`
}
