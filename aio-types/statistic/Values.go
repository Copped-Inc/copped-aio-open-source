package statistic

import "time"

type statistic struct {
	Id          string   `json:"id"`
	User        string   `json:"user"`
	Logs        []string `json:"logs"`
	Path        string   `json:"path"`
	Method      string   `json:"method"`
	RequestBody []byte   `json:"request_body"`

	Start    time.Time     `json:"start"`
	End      time.Time     `json:"end"`
	Duration time.Duration `json:"duration"`

	ClosingState state  `json:"closing_state"`
	Err          string `json:"err"`
	StatusCode   int    `json:"status_code"`
	ResponseBody []byte `json:"response_body"`
}

type log struct {
	Id      string    `json:"id"`
	State   state     `json:"state"`
	Ref     string    `json:"ref"`
	Time    time.Time `json:"time"`
	Content []any     `json:"content"`
}

type state int8

const (
	stateOk state = iota
	stateErr
	stateTimeout
)
