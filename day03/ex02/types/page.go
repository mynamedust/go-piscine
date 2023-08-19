package types

type JsonPage struct {
	Name         string  `json:"name,omitempty"`
	TotalCount   int     `json:"total,omitempty"`
	Places       []Place `json:"places,omitempty"`
	PreviousPage int     `json:"prev_page,omitempty"`
	NextPage     int     `json:"next_page,omitempty"`
	LastPage     int     `json:"last_page,omitempty"`
	Error        string  `json:"error,omitempty"`
}
