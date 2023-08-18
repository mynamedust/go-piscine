package types

type HtmlPage struct {
	TotalCount   int
	Places       []Place
	HasPrevious  bool
	HasNext      bool
	PreviousPage int
	NextPage     int
	LastPage     int
}
