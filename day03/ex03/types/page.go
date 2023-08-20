package types

type JsonPage struct {
	Name   string  `json:"name,omitempty"`
	Places []Place `json:"places,omitempty"`
	Error  string  `json:"error,omitempty"`
}
