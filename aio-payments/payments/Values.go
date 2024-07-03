package payments

type Payments struct {
	auth  string
	data  map[string]string
	Id    string `json:"id"`
	Store string `json:"store"`
	Url   string `json:"url"`
	Data  string `json:"data"`
	State state  `json:"state,omitempty"`
}

type state int

const (
	Created state = iota
	Accepted
	Declined
	Finalized
)
