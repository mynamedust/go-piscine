package types

type LocationType struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type Place struct {
	Name     string       `json:"name,omitempty"`
	Address  string       `json:"address,omitempty"`
	Phone    string       `json:"phone,omitempty"`
	Id       int          `json:"id,omitempty"`
	Location LocationType `json:"location,omitempty"`
}
