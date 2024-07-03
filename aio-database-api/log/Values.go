package log

import "time"

type Log struct {
	State    Op        `json:"state"`
	Date     time.Time `json:"date"`
	User     string    `json:"user,omitempty"`
	Instance string    `json:"instance,omitempty"`
	Message  string    `json:"message"`
}

type Op int

const (
	Info = iota
	Error
)
